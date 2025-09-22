package postgres

import (
	"context"
	"errors"
	"fmt"
	"math"
	"route256/cart/pkg/logger"
	"route256/loms/internal/domain"
	sqlcrepos "route256/loms/internal/infra/repository/postgres/sqlc/generated"

	"github.com/jackc/pgx/v5"
)

// Int64ToUint32 безопасно конвертирует int64 в uint32
func Int64ToUint32(num int64) (uint32, error) {
	if num < 0 {
		return 0, fmt.Errorf("invalid num value (num < 0): %d", num)
	}

	if num > math.MaxUint32 {
		return math.MaxUint32, fmt.Errorf("invalid num value (num > math.MaxUint32): %d", num)
	}

	return uint32(num), nil //nolint:gosec // G115: int64 -> uint32 checked manually
}

// NewStockRepository создает новый StockRepository.
func NewStockRepository(pool sqlcrepos.DBTX) *StockRepository {
	return &StockRepository{
		sqlcrepos.New(pool),
	}
}

// StockRepository предоставляет доступ к хранилищу запасов из postgres.
type StockRepository struct {
	querier sqlcrepos.Querier
}

// Upsert добавляет или обновляет запись о запасе в postgres.
func (sr *StockRepository) Upsert(ctx context.Context, stock *domain.Stock) error {
	err := sr.querier.AddStock(ctx, &sqlcrepos.AddStockParams{
		Sku:        stock.SkuID,
		TotalCount: int64(stock.TotalCount),
		Reserved:   int64(stock.Reserved),
	})
	if err != nil {
		return fmt.Errorf("querier.AddStock: %w", err)
	}

	return nil
}

// AddReserve резервирует товар по SKU в postgres.
func (sr *StockRepository) AddReserve(ctx context.Context, skuID int64, delta uint32) error {
	err := sr.querier.Reserve(ctx, &sqlcrepos.ReserveParams{
		Sku:      skuID,
		Reserved: int64(delta),
	})
	if err != nil {
		return fmt.Errorf("querier.Reserve: %w", err)
	}

	return nil
}

// RemoveReserve убирает резерв с товара по SKU в postgres.
func (sr *StockRepository) RemoveReserve(ctx context.Context, skuID int64, delta uint32) error {
	err := sr.querier.RemoveReserve(ctx, &sqlcrepos.RemoveReserveParams{
		Sku:      skuID,
		Reserved: int64(delta),
	})
	if err != nil {
		return fmt.Errorf("querier.RemoveReserve: %w", err)
	}

	return nil
}

// ReduceReserveAndTotal уменьшает резерв и общий запас товара по SKU в postgres.
func (sr *StockRepository) ReduceReserveAndTotal(ctx context.Context, skuID int64, delta uint32) error {
	err := sr.querier.ReduceTotalAndReserve(ctx, &sqlcrepos.ReduceTotalAndReserveParams{
		Sku:      skuID,
		Reserved: int64(delta),
	})
	if err != nil {
		return fmt.Errorf("querier.ReduceTotalAndReserve: %w", err)
	}

	return nil
}

// GetBySkuID возвращает запас по SKU из postgres.
func (sr *StockRepository) GetBySkuID(ctx context.Context, skuID int64) (*domain.Stock, error) {
	stockDB, err := sr.querier.GetStockBySKU(ctx, skuID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrItemStockNotExist
		}

		return nil, fmt.Errorf("querier.GetStockBySKU: %w", err)
	}

	totalCount, err := Int64ToUint32(stockDB.TotalCount)
	if err != nil {
		logger.Warning(fmt.Sprintf("Int64ToUint32 (TotalCount=%d): %s", stockDB.TotalCount, err.Error()))
	}

	reserved, err := Int64ToUint32(stockDB.Reserved)
	if err != nil {
		logger.Warning(fmt.Sprintf("Int64ToUint32 (Reserved=%d): %s", stockDB.Reserved, err.Error()))
	}

	return &domain.Stock{
		SkuID:      stockDB.Sku,
		TotalCount: totalCount,
		Reserved:   reserved,
	}, nil
}
