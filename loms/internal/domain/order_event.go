package domain

import "time"

// OrderEvent описывает событие по заказу для передачи во внешние системы.
type OrderEvent struct {
	OrderID int64  `json:"order_id"`
	Status  string `json:"status"`
	Moment  string `json:"moment"`
}

// OrderEventOutbox описывает событие заказа в outbox.
type OrderEventOutbox struct {
	ID      int64
	OrderID int64
	Status  string
	Moment  time.Time
}
