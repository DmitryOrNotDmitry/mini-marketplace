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
