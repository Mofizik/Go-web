package storage

import (
	"fmt"
	"order/internal/order/model"
	"sync"
)

type Storage struct {
	mu     sync.RWMutex
	orders map[string]*model.Order
}

func NewStorage() *Storage {
	return &Storage{
		orders: make(map[string]*model.Order),
	}
}

func (s *Storage) Save(order *model.Order) (*model.Order, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.orders[order.ID] = order
	return order, nil
}

func (s *Storage) FindByID(id string) (*model.Order, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	order, ok := s.orders[id]
	if !ok {
		return nil, fmt.Errorf("order with id %s not found", id)
	}
	return order, nil
}

func (s *Storage) Update(order *model.Order) (*model.Order, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.orders[order.ID]
	if !ok {
		return nil, fmt.Errorf("order with id %s not found", order.ID)
	}
	s.orders[order.ID] = order
	return order, nil
}

func (s *Storage) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.orders[id]
	if !ok {
		return fmt.Errorf("order with id %s not found", id)
	}
	delete(s.orders, id)
	return nil
}

func (s *Storage) List() ([]*model.Order, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]*model.Order, 0, len(s.orders))
	for _, o := range s.orders {
		result = append(result, o)
	}
	return result, nil
}
