package main

import (
	"google.golang.org/grpc"
	"net"
	"fmt"
	"log"
	"flag"
	pb "order/pkg/api/test"
	"order/pkg/api/v1"
	"order/internal/order"
)


var (
	port = flag.Int("port", 50051, "The server port")
)

func main () {
	storage := order.NewStorage()

	OrderService := order.NewOrderService(storage)
	server := v1.NewOrderHander(OrderService)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	pb.RegisterOrderServiceServer(s, server)

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}