package app

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"order/internal/order/api/handler"
	"order/internal/order/service"
	"order/internal/order/storage"
	pb "order/pkg/api/test"
	"order/pkg/config"
	"order/pkg/logger"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type App struct {
	grpcServer *grpc.Server
	lis        net.Listener
	log        *slog.Logger
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

	interceptor := func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {

		log.Info("info", "method", info.FullMethod, "req", req)

		resp, err = handler(ctx, req)

		log.Info("resp", resp, "err", err)

		return handler(ctx, req)
	}
	s := grpc.NewServer(grpc.UnaryInterceptor(interceptor))

	pb.RegisterOrderServiceServer(s, srv)
    reflection.Register(s)

	port := config.MustGet("GRPC_PORT")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		return nil, fmt.Errorf("app.New failed to listen: %w", err)
	}

	return &App{
		grpcServer: s,
		lis:        lis,
		log:        log,
	}, nil
}
func (a *App) Run() error {
	a.log.Info("Server listening", "addr", a.lis.Addr().String())

	if err := a.grpcServer.Serve(a.lis); err != nil {
		a.log.Error("failed to serve", "error", err)
		return err
	}
	return nil
}
