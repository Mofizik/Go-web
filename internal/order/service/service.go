package service

import (
	"order/internal/order/model"
	"order/internal/order/storage"
	"order/pkg/idgen"
)

type OrderService struct {
	stor *storage.Storage
}

func NewOrderService(stor *storage.Storage) *OrderService {
	return &OrderService{stor: stor}
}

func (s *OrderService) CreateOrder(item string, quantity int32) (string, error) {
	randID := idgen.GenerateId(16)
	order := &model.Order{
		ID:       randID,
		Item:     item,
		Quantity: quantity,
	}
	order, err := s.stor.Save(order)
	return order.ID, err
}

func (s *OrderService) GetOrder(id string) (*model.Order, error) {
	return s.stor.FindByID(id)
}

func (s *OrderService) UpdateOrder(order *model.Order) (*model.Order, error) {
	return s.stor.Update(order)
}

func (s *OrderService) DeleteOrder(id string) error {
	return s.stor.Delete(id)
}

func (s *OrderService) ListOrders() ([]*model.Order, error) {
	return s.stor.List()
}
