//go:build integration
// +build integration

package postgres

import (
	"context"
	"route256/loms/internal/domain"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostgresStockRepositoryIntegration(t *testing.T) {
	t.Parallel()

	pool, err := newTestPool(context.Background())
	require.NoError(t, err)

	stockRepository := NewStockRepository(pool)

	t.Run("insert stock and get", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		stock := &domain.Stock{
			SkuID:      1,
			TotalCount: 100,
			Reserved:   10,
		}

		err := stockRepository.Upsert(ctx, stock)
		assert.NoError(t, err)

		err = stockRepository.Upsert(ctx, stock)
		assert.NoError(t, err)

		actualStock, err := stockRepository.GetBySkuID(ctx, stock.SkuID)
		assert.NoError(t, err)

		deleteStock(ctx, pool, stock.SkuID)

		assert.Equal(t, stock.SkuID, actualStock.SkuID)
		assert.Equal(t, stock.TotalCount*2, actualStock.TotalCount)
		assert.Equal(t, stock.Reserved*2, actualStock.Reserved)
	})

	t.Run("add reverse for stock success", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		stock := &domain.Stock{
			SkuID:      2,
			TotalCount: 100,
			Reserved:   10,
		}

		err := stockRepository.Upsert(ctx, stock)
		assert.NoError(t, err)

		delta := uint32(50)
		err = stockRepository.AddReserve(ctx, stock.SkuID, delta)
		assert.NoError(t, err)

		actualStock, err := stockRepository.GetBySkuID(ctx, stock.SkuID)
		assert.NoError(t, err)

		deleteStock(ctx, pool, stock.SkuID)

		assert.Equal(t, stock.SkuID, actualStock.SkuID)
		assert.Equal(t, stock.TotalCount, actualStock.TotalCount)
		assert.Equal(t, stock.Reserved+delta, actualStock.Reserved)
	})

	t.Run("remove reverse for stock success", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		stock := &domain.Stock{
			SkuID:      3,
			TotalCount: 100,
			Reserved:   10,
		}

		err := stockRepository.Upsert(ctx, stock)
		assert.NoError(t, err)

		delta := uint32(10)
		err = stockRepository.RemoveReserve(ctx, stock.SkuID, delta)
		assert.NoError(t, err)

		actualStock, err := stockRepository.GetBySkuID(ctx, stock.SkuID)
		assert.NoError(t, err)

		deleteStock(ctx, pool, stock.SkuID)

		assert.Equal(t, stock.SkuID, actualStock.SkuID)
		assert.Equal(t, stock.TotalCount, actualStock.TotalCount)
		assert.Equal(t, stock.Reserved-delta, actualStock.Reserved)
	})

	t.Run("reduce reverse and total for stock success", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		stock := &domain.Stock{
			SkuID:      4,
			TotalCount: 100,
			Reserved:   10,
		}

		err := stockRepository.Upsert(ctx, stock)
		assert.NoError(t, err)

		delta := uint32(10)
		err = stockRepository.ReduceReserveAndTotal(ctx, stock.SkuID, delta)
		assert.NoError(t, err)

		actualStock, err := stockRepository.GetBySkuID(ctx, stock.SkuID)
		assert.NoError(t, err)

		deleteStock(ctx, pool, stock.SkuID)

		assert.Equal(t, stock.SkuID, actualStock.SkuID)
		assert.Equal(t, stock.TotalCount-delta, actualStock.TotalCount)
		assert.Equal(t, stock.Reserved-delta, actualStock.Reserved)
	})

}

func deleteStock(ctx context.Context, pool *pgxpool.Pool, sku int64) {
	pool.Exec(ctx, "delete from stocks where sku = $1", sku)
}
