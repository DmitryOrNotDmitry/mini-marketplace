package repository

import (
	"context"
	"route256/loms/internal/domain"
	"sync"
)

type OrderStorage = map[int64]*domain.Order

type IDGenerator interface {
	NextID() int64
}

type OrderRepositoryInMemory struct {
	generatorID IDGenerator
	storage     OrderStorage
	mx          sync.RWMutex
}

// NewInMemoryOrderRepository создает новый репозиторий заказов с in-memory хранилищем.
func NewInMemoryOrderRepository(generatorID IDGenerator, cap int) *OrderRepositoryInMemory {
	return &OrderRepositoryInMemory{
		generatorID: generatorID,
		storage:     make(OrderStorage, cap),
	}
}

// Insert добавляет новый заказ в хранилище и возвращает его идентификатор.
func (or *OrderRepositoryInMemory) Insert(_ context.Context, order *domain.Order) (int64, error) {
	or.mx.Lock()
	defer or.mx.Unlock()

	order.OrderID = or.generatorID.NextID()
	or.storage[order.OrderID] = order

	return order.OrderID, nil
}

// GetByID возвращает заказ по идентификатору, если он существует.
func (or *OrderRepositoryInMemory) GetByID(_ context.Context, orderID int64) (*domain.Order, error) {
	or.mx.Lock()
	defer or.mx.Unlock()

	order, ok := or.storage[orderID]
	if !ok {
		return nil, domain.ErrOrderNotExist
	}

	return order, nil
}

// UpdateStatus обновляет статус заказа по идентификатору.
func (or *OrderRepositoryInMemory) UpdateStatus(_ context.Context, orderID int64, newStatus domain.Status) error {
	or.mx.Lock()
	defer or.mx.Unlock()

	order, ok := or.storage[orderID]
	if !ok {
		return domain.ErrOrderNotExist
	}

	order.Status = newStatus

	return nil
}
