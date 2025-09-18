package service

import (
	"context"
	"fmt"
	"route256/loms/internal/domain"
)

// StockRepoFactory создает экземпляры репозитория запасов.
type StockRepoFactory interface {
	// CreateStock создает новый репозиторий запасов.
	CreateStock(ctx context.Context, operationType OperationType) StockRepository
}

// StockService реализует бизнес-логику управления запасами товаров.
type StockService struct {
	repositoryFactory StockRepoFactory
	txManager         TxManager
}

// NewStockService создает новый сервис управления запасами.
func NewStockService(repositoryFactory StockRepoFactory, txManager TxManager) *StockService {
	return &StockService{
		repositoryFactory: repositoryFactory,
		txManager:         txManager,
	}
}

// Create добавляет или обновляет запись о запасе.
func (ss *StockService) Create(ctx context.Context, stock *domain.Stock) error {
	if stock.Reserved > stock.TotalCount {
		return domain.ErrItemStockNotValid
	}

	stockRepository := ss.repositoryFactory.CreateStock(ctx, Write)
	err := stockRepository.Upsert(ctx, stock)
	if err != nil {
		return fmt.Errorf("stockRepository.Upsert: %w", err)
	}

	return nil
}

// GetAvailableCount возвращает количество доступного товара по SKU.
func (ss *StockService) GetAvailableCount(ctx context.Context, skuID int64) (uint32, error) {
	stockRepository := ss.repositoryFactory.CreateStock(ctx, Read)
	stock, err := stockRepository.GetBySkuID(ctx, skuID)
	if err != nil {
		return 0, fmt.Errorf("stockRepository.GetBySkuID: %w", err)
	}

	return stock.TotalCount - stock.Reserved, nil
}

// ReserveFor резервирует товары под заказ.
func (ss *StockService) ReserveFor(ctx context.Context, order *domain.Order) error {
	err := ss.txManager.WithRepeatableRead(ctx, Write, func(ctx context.Context) error {
		stockRepository := ss.repositoryFactory.CreateStock(ctx, FromTx)
		for _, item := range order.Items {
			stock, err := stockRepository.GetBySkuID(ctx, item.SkuID)
			if err != nil {
				return fmt.Errorf("stockRepository.GetBySkuID: %w", err)
			}

			if stock.TotalCount-stock.Reserved < item.Count {
				return domain.ErrCanNotReserveItem
			}

			err = stockRepository.AddReserve(ctx, item.SkuID, item.Count)
			if err != nil {
				return fmt.Errorf("stockRepository.AddReserve: %w", err)
			}
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("txManager.WithTransaction: %w", err)
	}

	return nil
}

// CancelReserveFor отменяет резервирование товаров по заказу.
func (ss *StockService) CancelReserveFor(ctx context.Context, order *domain.Order) error {
	err := ss.txManager.WithTransaction(ctx, Write, func(ctx context.Context) error {
		stockRepository := ss.repositoryFactory.CreateStock(ctx, FromTx)
		for _, item := range order.Items {
			err := stockRepository.RemoveReserve(ctx, item.SkuID, item.Count)
			if err != nil {
				return fmt.Errorf("stockRepository.RemoveReserve: %w", err)
			}
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("txManager.WithTransaction: %w", err)
	}

	return nil
}

// ConfirmReserveFor подтверждает резервирование и уменьшает общий запас.
func (ss *StockService) ConfirmReserveFor(ctx context.Context, order *domain.Order) error {
	err := ss.txManager.WithTransaction(ctx, Write, func(ctx context.Context) error {
		stockRepository := ss.repositoryFactory.CreateStock(ctx, FromTx)
		for _, item := range order.Items {
			err := stockRepository.ReduceReserveAndTotal(ctx, item.SkuID, item.Count)
			if err != nil {
				return fmt.Errorf("stockRepository.ReduceReserveAndTotal: %w", err)
			}
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("txManager.WithTransaction: %w", err)
	}

	return nil
}
