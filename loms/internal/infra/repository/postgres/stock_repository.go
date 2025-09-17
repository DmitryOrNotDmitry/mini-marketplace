package postgres

import (
	"context"
	"fmt"
	"route256/loms/internal/domain"
	repo_sqlc "route256/loms/internal/infra/repository/postgres/sqlc/generated"
)

// NewStockRepository создает новый StockRepository.
func NewStockRepository(pool repo_sqlc.DBTX) *StockRepository {
	return &StockRepository{
		repo_sqlc.New(pool),
	}
}

// StockRepository предоставляет доступ к хранилищу запасов из postgres.
type StockRepository struct {
	querier repo_sqlc.Querier
}

// Upsert добавляет или обновляет запись о запасе в postgres.
func (sr *StockRepository) Upsert(ctx context.Context, stock *domain.Stock) error {
	err := sr.querier.AddStock(ctx, &repo_sqlc.AddStockParams{
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
	err := sr.querier.Reserve(ctx, &repo_sqlc.ReserveParams{
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
	err := sr.querier.RemoveReserve(ctx, &repo_sqlc.RemoveReserveParams{
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
	err := sr.querier.ReduceTotalAndReserve(ctx, &repo_sqlc.ReduceTotalAndReserveParams{
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
		return nil, fmt.Errorf("querier.AddStock: %w", err)
	}

	return &domain.Stock{
		SkuID:      stockDB.Sku,
		TotalCount: uint32(stockDB.TotalCount),
		Reserved:   uint32(stockDB.Reserved),
	}, nil
}
