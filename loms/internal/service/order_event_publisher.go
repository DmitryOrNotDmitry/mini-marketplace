package service

import (
	"context"
	"encoding/json"
	"fmt"
	"route256/cart/pkg/logger"
	"route256/loms/internal/domain"
	"time"
)

type publisher interface {
	Send(key string, value []byte) error
}

type orderEventRepoFactory interface {
	CreateOrderEvent(ctx context.Context, operationType OperationType) OrderEventRepository
}

// OrderEventPublisher отвечает за публикацию событий заказов.
type OrderEventPublisher struct {
	pub               publisher
	repositoryFactory orderEventRepoFactory
	batchSize         int32
	period            time.Duration
}

// NewOrderEventPublisher создает новый экземпляр OrderEventPublisher.
func NewOrderEventPublisher(publisher publisher, repositoryFactory orderEventRepoFactory, batchSize int32, period time.Duration) *OrderEventPublisher {
	return &OrderEventPublisher{
		pub:               publisher,
		repositoryFactory: repositoryFactory,
		batchSize:         batchSize,
		period:            period,
	}
}

// Start запускает периодическую отправку событий заказов.
func (o *OrderEventPublisher) Start(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(o.period)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				err := o.sendEvents(ctx)
				if err != nil {
					logger.Warnw("error at sendEvents()", "err", err)
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (o *OrderEventPublisher) sendEvents(ctx context.Context) error {
	readOrderEventRepo := o.repositoryFactory.CreateOrderEvent(ctx, Read)
	events, err := readOrderEventRepo.GetUnprocessedEventsLimit(ctx, o.batchSize)
	if err != nil {
		return err
	}

	writeOrderEventRepo := o.repositoryFactory.CreateOrderEvent(ctx, Write)

	erroredOrders := make(map[int64]struct{})
	for _, event := range events {
		if _, ok := erroredOrders[event.OrderID]; ok {
			continue
		}

		msg := &domain.OrderEvent{
			OrderID: event.OrderID,
			Status:  event.Status,
			Moment:  event.Moment.Format(time.RFC3339),
		}

		msgBytes, err := json.Marshal(msg)
		if err != nil {
			erroredOrders[msg.OrderID] = struct{}{}
			continue
		}

		err = o.pub.Send(messageKey(msg), msgBytes)
		if err != nil {
			erroredOrders[msg.OrderID] = struct{}{}
			continue
		}

		err = writeOrderEventRepo.UpdateEventStatus(ctx, event.ID, domain.Complete)
		if err != nil {
			logger.WarnwCtx(ctx, fmt.Sprintf("Сообщение о статусе заказа с id=%d успешно отправлено в kafka, но статус в outbox не обновлен", event.ID), "err", err)
		}
	}

	return nil
}

func messageKey(msg *domain.OrderEvent) string {
	return fmt.Sprintf("%d", msg.OrderID)
}
