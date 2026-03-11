package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"order/internal/order/api/handler"
	"order/internal/order/service"
	"order/internal/order/storage"
	pb "order/pkg/api/test"
	"order/pkg/closer"
	"order/pkg/config"
	"order/pkg/grpcmw"
	"order/pkg/httpmw"
	"order/pkg/logger"
	"os"
	"os/signal"

	"syscall"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type App struct {
	grpcServer *grpc.Server
	httpServer *http.Server
	lis        net.Listener
	log        *slog.Logger
	closer	   *closer.Closer
}

func New(ctx context.Context) (*App, error) {

	//1. load env
	if err := config.LoadDotEnv("internal/order/config/.env"); err != nil {
		return nil, fmt.Errorf("app.New: %w", err)
	}
	env := config.Get("APP_ENV", "local")

	// 2. setup logger
	logger.Setup(env)
	log := logger.With("service", "order")

	// 3. create grpc server

	stor := storage.NewStorage()
	svc := service.NewOrderService(stor)
	srv := handler.NewOrderHandler(svc)

	s := grpc.NewServer(grpc.UnaryInterceptor(grpcmw.UnaryServerLoggingInterceptor(log)))

	// 4. create http server
	gwMux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	
	

	pb.RegisterOrderServiceServer(s, srv)
    reflection.Register(s)

	grpcPort := config.MustGet("GRPC_PORT")
	pb.RegisterOrderServiceHandlerFromEndpoint(ctx, gwMux, fmt.Sprintf("localhost:%s", grpcPort), opts) // Register http handler to grpc gateway

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		return nil, fmt.Errorf("app.New failed to listen: %w", err)
	}
	
	httpPort := config.MustGet("HTTP_PORT")
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", httpPort),
		Handler: httpmw.LoggingMiddleware(log)(gwMux),
	}

	shutdownCloser := closer.New(log)
	
	return &App{
		grpcServer: s,
		httpServer: httpServer,
		lis:        lis,
		log:        log,
		closer:	   shutdownCloser,
	}, nil
}
func (a *App) Run() error {
    a.closer.AddFunc("grpc listener", func() {
        _ = a.lis.Close()
    })
    a.closer.Add("close grpc server", func(ctx context.Context) error {
        done := make(chan struct{})
        go func() {
            a.grpcServer.GracefulStop()
            close(done)
        }()
        select {
        case <-done:
            return nil
        case <-ctx.Done():
            a.grpcServer.Stop()
            <-done
            return ctx.Err()
        }
    })
	
	a.closer.Add("http server", func(ctx context.Context) error {
		return a.httpServer.Shutdown(ctx)
	})

    sigCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
    defer stop()

    go func() {
        <-sigCtx.Done()
        a.log.Info("shutdown signal received")
        stop()

        shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()

        if err := a.closer.Close(shutdownCtx); err != nil && !errors.Is(err, context.DeadlineExceeded) {
            a.log.Error("graceful shutdown error", "error", err)
        }
        a.log.Info("app.Run shutdown complete")
		fmt.Println("Server Stopped")
    }()
	
	a.log.Info("http server listening", "addr", a.httpServer.Addr)
	go a.httpServer.ListenAndServe()

    a.log.Info("gRPC server listening", "addr", a.lis.Addr().String())
    if err := a.grpcServer.Serve(a.lis); err != nil {
        if sigCtx.Err() != nil {
            return nil
        }
        return fmt.Errorf("app.Run: %w", err)
    }

	


    return nil
}
