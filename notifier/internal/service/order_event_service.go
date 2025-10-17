package service

import (
	"route256/cart/pkg/logger"
	"route256/notifier/internal/domain"
)

// OrderEventService реализует обработку входящих событий заказа.
type OrderEventService struct {
	msgIn <-chan *domain.OrderEvent
}

// NewOrderEventService создает новый экземпляр OrderEventConsumer.
func NewOrderEventService(msgIn <-chan *domain.OrderEvent) *OrderEventService {
	return &OrderEventService{
		msgIn: msgIn,
	}
}

// Start запускает обработку событий заказа.
func (o *OrderEventService) Start() {
	go func() {
		for msg := range o.msgIn {
			logger.Infow("Событие заказа изменилось", "order_id", msg.OrderID, "status", msg.Status, "moment", msg.Moment)
		}
	}()
}
