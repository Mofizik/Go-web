package order


type StorageInterface interface {
	Save(order *Order)
	FindByID(id string)
	Update(order *Order)
	Delete(id string)
	List()
}