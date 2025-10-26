package kafka

// OrderEventKafka описывает событие по заказу для kafka.
type OrderEventKafka struct {
	OrderID int64  `json:"order_id"`
	Status  string `json:"status"`
	Moment  string `json:"moment"`
}
