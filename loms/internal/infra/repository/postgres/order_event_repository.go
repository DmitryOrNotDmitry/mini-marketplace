package postgres

import (
	"context"
	"fmt"
	"route256/loms/internal/domain"
	sqlcrepos "route256/loms/internal/infra/repository/postgres/sqlc/generated"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// NewOrderEventRepository создает новый OrderEventRepository.
func NewOrderEventRepository(pool sqlcrepos.DBTX) *OrderEventRepository {
	return &OrderEventRepository{
		sqlcrepos.New(pool),
	}
}

// OrderEventRepository предоставляет доступ к хранилищу событий о смене статуса заказа (outbox) из postgres.
type OrderEventRepository struct {
	querier sqlcrepos.Querier
}

// Insert добавляет новое событие об изменении статуса заказа в outbox.
func (oe *OrderEventRepository) Insert(ctx context.Context, order *domain.Order) error {
	err := oe.querier.InsertOrderEvent(ctx, &sqlcrepos.InsertOrderEventParams{
		OrderID:     &order.OrderID,
		Status:      string(order.Status),
		Moment:      now(),
		EventStatus: string(domain.EventNew),
	})
	if err != nil {
		return fmt.Errorf("querier.InsertOrderEvent: %w", err)
	}

	return nil
}

func now() pgtype.Timestamp {
	return pgtype.Timestamp{
		Time:  time.Now(),
		Valid: true,
	}
}

// GetUnprocessedEventsLimit возращает первые limit событий, которые еще не обработаны.
func (oe *OrderEventRepository) GetUnprocessedEventsLimit(ctx context.Context, limit int) ([]*domain.OrderEventOutbox, error) {
	rows, err := oe.querier.GetUnprocessedEventsLimit(ctx, int32(limit))
	if err != nil {
		return nil, fmt.Errorf("querier.GetUnprocessedEventsLimit: %w", err)
	}

	res := make([]*domain.OrderEventOutbox, 0, len(rows))
	for _, row := range rows {
		res = append(res, &domain.OrderEventOutbox{
			ID:      row.ID,
			OrderID: *row.OrderID,
			Status:  row.Status,
			Moment:  row.Moment.Time,
		})
	}

	return res, nil
}

// UpdateEventStatus обновляет статус события.
func (oe *OrderEventRepository) UpdateEventStatus(ctx context.Context, eventID int64, newStatus domain.EventStatus) error {
	err := oe.querier.UpdateEventStatus(ctx, &sqlcrepos.UpdateEventStatusParams{
		ID:          eventID,
		EventStatus: string(newStatus),
	})
	if err != nil {
		return fmt.Errorf("querier.UpdateEventStatus: %w", err)
	}

	return nil
}
