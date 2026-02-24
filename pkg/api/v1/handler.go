package v1

import (
	"context"
	"order/internal/order"
	pb "order/pkg/api/test"
)
type OrderHander struct {
	pb.UnimplementedOrderServiceServer
	OrderService *order.OrderService
}

func NewOrderHander(os *order.OrderService) *OrderHander {
	return &OrderHander{OrderService: os}
}

func convertOrderStruct(o *order.Order) (*pb.Order){
	return &pb.Order{
		Id: o.ID,
		Item: o.Item,
		Quantity: o.Quantity,
	}
}



func (h *OrderHander) CreateOrder (ctx context.Context, r *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	res, err := h.OrderService.CreateOrder(r.Item, r.Quantity)
	if err != nil {
		return nil, err 
	}
	result := &pb.CreateOrderResponse{
		Id: res,
	}
	return result, err
}

func (h *OrderHander) GetOrder (ctx context.Context, r *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
	res, err := h.OrderService.GetOrder(r.Id)
	if err != nil {
		return nil, err 
	}
	order := convertOrderStruct(res)

	result := &pb.GetOrderResponse{
		Order: order,
	}
	return result, err
}

func (h *OrderHander) UpdateOrder (ctx context.Context, r *pb.UpdateOrderRequest) (*pb.UpdateOrderResponse, error) {
	req := &order.Order{
		ID: r.Id,
		Item: r.Item,
		Quantity: r.Quantity,
	}
	res, err := h.OrderService.UpdateOrder(req)
	if err != nil {
		return nil, err 
	}
	order := convertOrderStruct(res)

	result := &pb.UpdateOrderResponse{
		Order: order,
	}
	return result, err
}


func (h *OrderHander) DeleteOrder (ctx context.Context, r *pb.DeleteOrderRequest) (*pb.DeleteOrderResponse, error) {
	err := h.OrderService.DeleteOrder(r.Id)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteOrderResponse{ Success: true}, err 
}


func (h *OrderHander) ListOrders (ctx context.Context, r *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	res, err := h.OrderService.ListOrders()
	if err != nil {
		return nil, err 
	}

	response := make([]*pb.Order, 0, len(res))
    for _, o := range res {
		response = append(response, convertOrderStruct(o))
	}
	result := &pb.ListOrdersResponse{
		Orders: response,
	}
	return result, err
}