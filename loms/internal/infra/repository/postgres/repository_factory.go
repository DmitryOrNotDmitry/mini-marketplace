package postgres

import (
	"context"
	"route256/loms/internal/service"

	"github.com/jackc/pgx/v5/pgxpool"
)

// PoolManager определяет интерфейс для получения пулов соединений.
type PoolManager interface {
	Readable() *pgxpool.Pool
	Writable() *pgxpool.Pool
}

// RepositoryFactory создает репозитории с нужным пулом соединений.
type RepositoryFactory struct {
	poolManager PoolManager
}

// NewRepositoryFactory создает новый RepositoryFactory.
func NewRepositoryFactory(poolManager PoolManager) *RepositoryFactory {
	return &RepositoryFactory{
		poolManager: poolManager,
	}
}

func (rf *RepositoryFactory) getPool(operationType service.OperationType) *pgxpool.Pool {
	if operationType == service.Read {
		return rf.poolManager.Readable()
	}

	return rf.poolManager.Writable()
}

// CreateStock создает StockRepository с нужным пулом или транзакцией.
func (rf *RepositoryFactory) CreateStock(ctx context.Context, operationType service.OperationType) service.StockRepository {
	if tx, ok := TxFromCtx(ctx); ok {
		return NewStockRepository(tx)
	}
	return NewStockRepository(rf.getPool(operationType))
}

// CreateOrder создает OrderRepository с нужным пулом или транзакцией.
func (rf *RepositoryFactory) CreateOrder(ctx context.Context, operationType service.OperationType) service.OrderRepository {
	if tx, ok := TxFromCtx(ctx); ok {
		return NewOrderRepository(tx)
	}
	return NewOrderRepository(rf.getPool(operationType))
}

// CreateOrderEvent создает OrderEventRepository с нужным пулом или транзакцией.
func (rf *RepositoryFactory) CreateOrderEvent(ctx context.Context, operationType service.OperationType) service.OrderEventRepository {
	if tx, ok := TxFromCtx(ctx); ok {
		return NewOrderEventRepository(tx)
	}
	return NewOrderEventRepository(rf.getPool(operationType))
}
