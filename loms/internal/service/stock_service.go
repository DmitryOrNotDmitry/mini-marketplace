package service

import (
	"context"
	"fmt"
	"route256/cart/pkg/logger"
	"route256/loms/internal/domain"
)

// StockRepository описывает методы работы с запасами товаров в хранилище.
type StockRepository interface {
	// Upsert добавляет или обновляет запись о запасе.
	Upsert(ctx context.Context, stock *domain.Stock) error
	// AddReserve резервирует товар по SKU.
	AddReserve(ctx context.Context, skuID int64, delta uint32) error
	// RemoveReserve убирает резерв с товара по SKU.
	RemoveReserve(ctx context.Context, skuID int64, delta uint32) error
	// ReduceReserveAndTotal уменьшает резерв и общий запас товара по SKU.
	ReduceReserveAndTotal(ctx context.Context, skuID int64, delta uint32) error
	// GetBySkuID возвращает информацию о запасе по SKU.
	GetBySkuID(ctx context.Context, skuID int64) (*domain.Stock, error)
}

// StockService реализует бизнес-логику управления запасами товаров.
type StockService struct {
	stockRepository StockRepository
}

// NewStockService создает новый сервис управления запасами.
func NewStockService(stockRepository StockRepository) *StockService {
	return &StockService{
		stockRepository: stockRepository,
	}
}

// Create добавляет или обновляет запись о запасе.
func (ss *StockService) Create(ctx context.Context, stock *domain.Stock) error {
	if stock.Reserved > stock.TotalCount {
		return domain.ErrItemStockNotValid
	}

	err := ss.stockRepository.Upsert(ctx, stock)
	if err != nil {
		return fmt.Errorf("stockRepository.Upsert: %w", err)
	}

	return nil
}

// GetAvailableCount возвращает количество доступного товара по SKU.
func (ss *StockService) GetAvailableCount(ctx context.Context, skuID int64) (uint32, error) {
	stock, err := ss.stockRepository.GetBySkuID(ctx, skuID)
	if err != nil {
		return 0, fmt.Errorf("stockRepository.GetBySkuID: %w", err)
	}

	return stock.TotalCount - stock.Reserved, nil
}

// ReserveFor резервирует товары под заказ.
func (ss *StockService) ReserveFor(ctx context.Context, order *domain.Order) error {
	var err error
	reservedItems := make([]*domain.OrderItem, 0, len(order.Items))

	for _, item := range order.Items {
		err = ss.stockRepository.AddReserve(ctx, item.SkuID, item.Count)
		if err != nil {
			break
		}

		reservedItems = append(reservedItems, item)
	}

	if err != nil {
		ss.rollbackReserve(ctx, reservedItems)
		return err
	}

	return nil
}

func (ss *StockService) rollbackReserve(ctx context.Context, reservedItems []*domain.OrderItem) {
	for _, item := range reservedItems {
		err := ss.stockRepository.RemoveReserve(ctx, item.SkuID, item.Count)
		if err != nil {
			logger.Error(fmt.Sprintf("stockRepository.RemoveReserve: %s", err.Error()))
		}
	}
}

// CancelReserveFor отменяет резервирование товаров по заказу.
func (ss *StockService) CancelReserveFor(ctx context.Context, order *domain.Order) error {
	var err error
	for _, item := range order.Items {
		iErr := ss.stockRepository.RemoveReserve(ctx, item.SkuID, item.Count)
		if iErr != nil {
			err = fmt.Errorf("stockRepository.RemoveReserve: %w", err)
		}
	}

	return err
}

// ConfirmReserveFor подтверждает резервирование и уменьшает общий запас.
func (ss *StockService) ConfirmReserveFor(ctx context.Context, order *domain.Order) error {
	var err error
	for _, item := range order.Items {
		iErr := ss.stockRepository.ReduceReserveAndTotal(ctx, item.SkuID, item.Count)
		if iErr != nil {
			err = fmt.Errorf("stockRepository.ReduceReserveAndTotal: %w", err)
		}
	}

	return err
}
