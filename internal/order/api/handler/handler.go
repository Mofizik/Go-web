package handler

import (
	"context"
	"order/internal/order/model"
	pb "order/pkg/api/test"
	"google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

type OrderServiceInterface interface {
	CreateOrder(ctx context.Context, item string, quantity int32) (string, error)
	GetOrder(ctx context.Context, id string) (*model.Order, error)
	UpdateOrder(ctx context.Context, o *model.Order) (*model.Order, error)
	DeleteOrder(ctx context.Context, id string) error
	ListOrders(ctx context.Context) ([]*model.Order, error)
}

type OrderHandler struct {
	pb.UnimplementedOrderServiceServer
	svc OrderServiceInterface
}

func NewOrderHandler(svc OrderServiceInterface) *OrderHandler {
	return &OrderHandler{svc: svc}
}

func convertOrderStruct(o *model.Order) *pb.Order {
	return &pb.Order{
		Id:       o.ID,
		Item:     o.Item,
		Quantity: o.Quantity,
	}
}

func (h *OrderHandler) CreateOrder(ctx context.Context, r *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	id, err := h.svc.CreateOrder(ctx, r.Item, r.Quantity)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create order: %v", err)
	}
	return &pb.CreateOrderResponse{Id: id}, nil
}

func (h *OrderHandler) GetOrder(ctx context.Context, r *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
    o, err := h.svc.GetOrder(ctx, r.Id)
    if err != nil {
        return nil, status.Errorf(codes.NotFound, "order not found: %s", r.Id)
    }
    return &pb.GetOrderResponse{Order: convertOrderStruct(o)}, nil
}


func (h *OrderHandler) UpdateOrder(ctx context.Context, r *pb.UpdateOrderRequest) (*pb.UpdateOrderResponse, error) {
    o, err := h.svc.UpdateOrder(ctx, &model.Order{
        ID:       r.Id,
        Item:     r.Item,
        Quantity: r.Quantity,
    })
    if err != nil {
        return nil, status.Errorf(codes.NotFound, "order not found: %s", r.Id)
    }
    return &pb.UpdateOrderResponse{Order: convertOrderStruct(o)}, nil
}

func (h *OrderHandler) DeleteOrder(ctx context.Context, r *pb.DeleteOrderRequest) (*pb.DeleteOrderResponse, error) {
	err := h.svc.DeleteOrder(ctx, r.Id)
    if err != nil {
        return nil, status.Errorf(codes.NotFound, "order not found: %s", r.Id)
    }
    return &pb.DeleteOrderResponse{Success: true}, nil
}

func (h *OrderHandler) ListOrders(ctx context.Context, r *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	orders, err := h.svc.ListOrders(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list orders: %v", err)
	}
	response := make([]*pb.Order, 0, len(orders))
	for _, o := range orders {
		response = append(response, convertOrderStruct(o))
	}
	return &pb.ListOrdersResponse{Orders: response}, nil
}
