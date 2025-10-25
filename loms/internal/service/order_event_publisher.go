package service

import (
	"context"
	"fmt"
	"route256/cart/pkg/logger"
	"route256/loms/internal/domain"
	"time"
)

type publisher interface {
	Send(key string, value *domain.OrderEvent) error
}

type orderEventRepoFactory interface {
	CreateOrderEvent(ctx context.Context, operationType OperationType) OrderEventRepository
}

// OrderEventPublisher отвечает за публикацию событий заказов.
type OrderEventPublisher struct {
	pub               publisher
	txManager         TxManager
	repositoryFactory orderEventRepoFactory
	batchSize         int32
	period            time.Duration
}

// NewOrderEventPublisher создает новый экземпляр OrderEventPublisher.
func NewOrderEventPublisher(publisher publisher, txManager TxManager, repositoryFactory orderEventRepoFactory, batchSize int32, period time.Duration) *OrderEventPublisher {
	return &OrderEventPublisher{
		pub:               publisher,
		txManager:         txManager,
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

	statuses := make(map[int64]domain.EventStatus, len(events))
	erroredOrders := make(map[int64]struct{}, len(events))

	for _, event := range events {
		if _, ok := erroredOrders[event.OrderID]; ok {
			statuses[event.ID] = domain.Dead
			continue
		}

		msg := &domain.OrderEvent{
			OrderID: event.OrderID,
			Status:  event.OrderStatus,
			Moment:  event.Moment.Format(time.RFC3339),
		}

		innerErr := o.pub.Send(messageKey(msg), msg)
		if innerErr != nil {
			statuses[event.ID] = domain.Dead
			erroredOrders[msg.OrderID] = struct{}{}
			continue
		}

		statuses[event.ID] = domain.Complete
	}

	return o.updateEventsStatusesTx(ctx, events, statuses)
}

func (o *OrderEventPublisher) updateEventsStatusesTx(ctx context.Context, events []*domain.OrderEventOutbox, statuses map[int64]domain.EventStatus) error {
	completedIDs := make([]int64, 0, len(events))
	deadIDs := make([]int64, 0, len(events))

	for id, status := range statuses {
		if status == domain.Complete {
			completedIDs = append(completedIDs, id)
			continue
		}

		if status == domain.Dead {
			deadIDs = append(deadIDs, id)
			continue
		}
	}

	err := o.txManager.WithTransaction(ctx, Write, func(ctx context.Context) error {
		writeOrderEventRepo := o.repositoryFactory.CreateOrderEvent(ctx, FromTx)

		err := writeOrderEventRepo.UpdateEventStatusBatch(ctx, completedIDs, domain.Complete)
		if err != nil {
			return fmt.Errorf("не удалось обновить статусы событий заказов в outbox: %w", err)
		}

		err = writeOrderEventRepo.UpdateEventStatusBatch(ctx, deadIDs, domain.Dead)
		if err != nil {
			return fmt.Errorf("не удалось обновить статусы событий заказов в outbox: %w", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("txManager.WithTransaction: %w", err)
	}

	return nil
}

func messageKey(msg *domain.OrderEvent) string {
	return fmt.Sprintf("%d", msg.OrderID)
}
