package handler

import (
	"context"
	"order/internal/order/model"
	pb "order/pkg/api/test"
)

type OrderServiceInterface interface {
	CreateOrder(item string, quantity int32) (string, error)
	GetOrder(id string) (*model.Order, error)
	UpdateOrder(o *model.Order) (*model.Order, error)
	DeleteOrder(id string) error
	ListOrders() ([]*model.Order, error)
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
	id, err := h.svc.CreateOrder(r.Item, r.Quantity)
	if err != nil {
		return nil, err
	}
	return &pb.CreateOrderResponse{Id: id}, nil
}

func (h *OrderHandler) GetOrder(ctx context.Context, r *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
	o, err := h.svc.GetOrder(r.Id)
	if err != nil {
		return nil, err
	}
	return &pb.GetOrderResponse{Order: convertOrderStruct(o)}, nil
}

func (h *OrderHandler) UpdateOrder(ctx context.Context, r *pb.UpdateOrderRequest) (*pb.UpdateOrderResponse, error) {
	o, err := h.svc.UpdateOrder(&model.Order{
		ID:       r.Id,
		Item:     r.Item,
		Quantity: r.Quantity,
	})
	if err != nil {
		return nil, err
	}
	return &pb.UpdateOrderResponse{Order: convertOrderStruct(o)}, nil
}

func (h *OrderHandler) DeleteOrder(ctx context.Context, r *pb.DeleteOrderRequest) (*pb.DeleteOrderResponse, error) {
	if err := h.svc.DeleteOrder(r.Id); err != nil {
		return nil, err
	}
	return &pb.DeleteOrderResponse{Success: true}, nil
}

func (h *OrderHandler) ListOrders(ctx context.Context, r *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	orders, err := h.svc.ListOrders()
	if err != nil {
		return nil, err
	}
	response := make([]*pb.Order, 0, len(orders))
	for _, o := range orders {
		response = append(response, convertOrderStruct(o))
	}
	return &pb.ListOrdersResponse{Orders: response}, nil
}
