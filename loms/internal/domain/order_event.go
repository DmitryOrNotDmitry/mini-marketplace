package domain

import "time"

// OrderEvent описывает событие заказа.
type OrderEvent struct {
	OrderID int64
	Status  string
	Moment  string
}

// OrderEventOutbox описывает событие заказа в outbox.
type OrderEventOutbox struct {
	ID          int64
	OrderID     int64
	OrderStatus string
	Moment      time.Time
}
