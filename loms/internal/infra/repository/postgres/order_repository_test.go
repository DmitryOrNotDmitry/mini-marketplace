//go:build integration
// +build integration

package postgres

import (
	"context"
	"fmt"
	"route256/loms/internal/domain"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestPool(ctx context.Context) (*pgxpool.Pool, error) {
	dsn := "postgresql://loms-user:loms-password@localhost:5432/loms_db?sslmode=disable"

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.ParseConfig (dsn=%s): %w", dsn, err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.NewWithConfig (dsn=%s): %w", dsn, err)
	}

	return pool, nil
}

func TestPostgresOrderRepositoryIntegration(t *testing.T) {
	t.Parallel()

	pool, err := newTestPool(context.Background())
	require.NoError(t, err)

	orderRepository := NewOrderRepository(pool)

	t.Run("insert order and get", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		order := &domain.Order{
			UserID: 1,
			Items: []*domain.OrderItem{
				&domain.OrderItem{SkuID: 1, Count: 100},
				&domain.OrderItem{SkuID: 2, Count: 100},
				&domain.OrderItem{SkuID: 3, Count: 100},
			},
		}

		orderID, err := orderRepository.Insert(ctx, order)
		require.NoError(t, err)
		order.OrderID = orderID

		actualOrder, err := orderRepository.GetByIDOrderItemsBySKU(ctx, orderID)
		assert.NoError(t, err)

		deleteOrder(ctx, pool, orderID)

		assert.Equal(t, order, actualOrder)
	})

	t.Run("insert order and update status", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		order := &domain.Order{
			UserID: 1,
			Items:  []*domain.OrderItem{},
			Status: domain.AwaitingPayment,
		}

		orderID, err := orderRepository.Insert(ctx, order)
		require.NoError(t, err)

		newStatus := domain.Cancelled

		err = orderRepository.UpdateStatus(ctx, orderID, newStatus)
		assert.NoError(t, err)

		actualOrder, err := orderRepository.GetByIDOrderItemsBySKU(ctx, orderID)
		assert.NoError(t, err)

		deleteOrder(ctx, pool, orderID)

		assert.Equal(t, newStatus, actualOrder.Status)
	})

}

func deleteOrder(ctx context.Context, pool *pgxpool.Pool, orderID int64) {
	pool.Exec(ctx, "delete from orders where order_id = $1", orderID)
}
