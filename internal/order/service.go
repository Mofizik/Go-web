package order

import "order/utils"

type OrderService struct {
    stor *Storage
}


func NewOrderService(stor *Storage) *OrderService {
    return &OrderService{stor: stor}
}


func (s *OrderService) CreateOrder (item string, quantity int32) (string, error) {
    randID := utils.GenerateId(16)
    order := &Order{
        ID : randID,
        Item: item,
        Quantity: quantity,
    }
    order, err := s.stor.Save(order)
    return order.ID, err
}

func (s *OrderService) GetOrder (id string) (*Order, error){
    return s.stor.FindByID(id)
}

func (s *OrderService) UpdateOrder(order *Order) (*Order, error){
    return s.stor.Update(order)
}

func (s *OrderService) DeleteOrder(id string ) (error){
    return s.stor.Delete(id)
}

func (s *OrderService) ListOrders() ([]*Order, error){
    return s.stor.List()
}
