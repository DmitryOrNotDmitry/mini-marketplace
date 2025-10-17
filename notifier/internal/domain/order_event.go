package domain

// OrderEvent описывает событие по заказу.
type OrderEvent struct {
	OrderID int64  `json:"order_id"`
	Status  string `json:"status"`
	Moment  string `json:"moment"`
}
