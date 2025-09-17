package postgres

import (
	"context"
	"route256/loms/internal/service"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RepositoryFactory struct {
	pool *pgxpool.Pool
}

func NewRepositoryFactory(pool *pgxpool.Pool) *RepositoryFactory {
	return &RepositoryFactory{
		pool: pool,
	}
}

func (rf *RepositoryFactory) CreateStock(ctx context.Context) service.StockRepository {
	if tx, ok := TxFromCtx(ctx); ok {
		return NewStockRepository(tx)
	}
	return NewStockRepository(rf.pool)
}

func (rf *RepositoryFactory) CreateOrder(ctx context.Context) service.OrderRepository {
	if tx, ok := TxFromCtx(ctx); ok {
		return NewOrderRepository(tx)
	}
	return NewOrderRepository(rf.pool)
}
