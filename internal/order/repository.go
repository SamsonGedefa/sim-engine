package order



type Repository interface {
	StoreOrder(order Order) error
	RetrieveOrder(orderID int) (Order, error)
	UpdateOrder(order Order) error
}


