package app

import (
    "flag"
    "fmt"
    "log"
    "net"
	"order/internal/order/api/handler"
    "google.golang.org/grpc"
    "order/internal/order/storage"
    "order/internal/order/service"
    pb "order/pkg/api/test"
)

func Run() {
    port := flag.Int("port", 50051, "The server port")
    flag.Parse() // ← не забыть

    stor := storage.NewStorage()
    svc := service.NewOrderService(stor)
    srv := handler.NewOrderHandler(svc)

    lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

    s := grpc.NewServer()
    pb.RegisterOrderServiceServer(s, srv)

    log.Printf("server listening at %v", lis.Addr())
    if err := s.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}