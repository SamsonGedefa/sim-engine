package order

import (
	"time"
)

type (

	OrderType int

	OrderStatus int

	Order struct {
		OrderID   int
		Symbol    string
		Price     float64
		Quantity  int
		Timestamp time.Time
		Status    OrderStatus
		Type      OrderType
	}

	OrderHandler interface {
		HandleOrder(order Order)
	}
)

const (
	Submitted OrderStatus = iota
	PartiallyFilled
	Filled
	Canceled
)



type Service struct {
	repo Repository
}


func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

