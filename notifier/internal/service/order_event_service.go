package service

import (
	"route256/cart/pkg/logger"
	"route256/notifier/internal/domain"
)

// OrderEventService реализует обработку входящих событий заказа.
type OrderEventService struct {
}

// NewOrderEventService создает новый экземпляр OrderEventConsumer.
func NewOrderEventService() *OrderEventService {
	return &OrderEventService{}
}

// Process логгирует событие о заказе в лог
func (o *OrderEventService) Process(event *domain.OrderEvent) {
	logger.Infow("Событие заказа изменилось", "order_id", event.OrderID, "status", event.Status, "moment", event.Moment)
}
