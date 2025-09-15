package repository

import (
	"context"
	"route256/loms/internal/domain"
	"sync"
)

// Storage хранит запасы по ID.
type Storage = map[int64]*domain.Stock

// StockRepositoryInMemory хранит запасы и резервирование в in-memory хранилище.
type StockRepositoryInMemory struct {
	storage Storage
	mx      sync.RWMutex
}

// NewInMemoryStockRepository создает новый репозиторий стоков с in-memory хранилищем.
func NewInMemoryStockRepository(cap int) *StockRepositoryInMemory {
	return &StockRepositoryInMemory{
		storage: make(Storage, cap),
	}
}

// Upsert добавляет или обновляет запись о запасе.
func (sr *StockRepositoryInMemory) Upsert(_ context.Context, stock *domain.Stock) error {
	sr.mx.Lock()
	defer sr.mx.Unlock()

	existingStock, ok := sr.storage[stock.SkuID]
	if ok {
		existingStock.TotalCount += stock.TotalCount
		existingStock.Reserved += stock.Reserved
	} else {
		sr.storage[stock.SkuID] = stock
	}

	return nil
}

// AddReserve резервирует указанное количество товара.
func (sr *StockRepositoryInMemory) AddReserve(_ context.Context, skuID int64, delta uint32) error {
	sr.mx.Lock()
	defer sr.mx.Unlock()

	stock, ok := sr.storage[skuID]
	if !ok {
		return domain.ErrCanNotReserveItem
	}

	if stock.TotalCount-stock.Reserved < delta {
		return domain.ErrCanNotReserveItem
	}

	stock.Reserved += delta

	return nil
}

// RemoveReserve снимает резервирование с товара.
func (sr *StockRepositoryInMemory) RemoveReserve(_ context.Context, skuID int64, delta uint32) error {
	sr.mx.Lock()
	defer sr.mx.Unlock()

	stock, ok := sr.storage[skuID]
	if ok {
		stock.Reserved -= delta
	}

	return nil
}

// ReduceReserveAndTotal уменьшает резерв и общий запас товара.
func (sr *StockRepositoryInMemory) ReduceReserveAndTotal(_ context.Context, skuID int64, delta uint32) error {
	sr.mx.Lock()
	defer sr.mx.Unlock()

	stock, ok := sr.storage[skuID]
	if ok {
		stock.Reserved -= delta
		stock.TotalCount -= delta
	}

	return nil
}

// GetBySkuID возвращает запас по SKU.
func (sr *StockRepositoryInMemory) GetBySkuID(_ context.Context, skuID int64) (*domain.Stock, error) {
	sr.mx.RLock()
	defer sr.mx.RUnlock()

	stock, ok := sr.storage[skuID]
	if ok {
		return stock, nil
	}

	return nil, domain.ErrItemStockNotExist
}
