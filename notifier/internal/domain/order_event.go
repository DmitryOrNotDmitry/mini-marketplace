package domain

// OrderEvent описывает событие по заказу.
type OrderEvent struct {
	OrderID int64
	Status  string
	Moment  string
}
