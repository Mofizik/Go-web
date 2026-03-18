package service

import (
	"context"
	"order/internal/order/model"
	"order/pkg/idgen"
)



type OrderRepository interface {
    Save(ctx context.Context, order *model.Order) (*model.Order, error)
    FindByID(ctx context.Context, id string) (*model.Order, error)
    Update(ctx context.Context, order *model.Order) (*model.Order, error)
    Delete(ctx context.Context, id string) error
    List(ctx context.Context) ([]*model.Order, error)
}

type OrderService struct {
	repo OrderRepository
}

func NewOrderService(repo OrderRepository) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) CreateOrder(ctx context.Context, item string, quantity int32) (string, error) {
	randID := idgen.GenerateId(16)
	order := &model.Order{
		ID:       randID,
		Item:     item,
		Quantity: quantity,
	}
	order, err := s.repo.Save(ctx, order)
	if err != nil {
		return "", err
	}
	return order.ID, err
}

func (s *OrderService) GetOrder(ctx context.Context, id string) (*model.Order, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *OrderService) UpdateOrder(ctx context.Context, order *model.Order) (*model.Order, error) {
	return s.repo.Update(ctx, order)
}

func (s *OrderService) DeleteOrder(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *OrderService) ListOrders(ctx context.Context) ([]*model.Order, error) {
	return s.repo.List(ctx)
}
