package repository

import (
	"context"
	"route256/loms/internal/domain"
	"sort"
	"sync"
)

// OrderStorage хранит заказы по ID.
type OrderStorage = map[int64]*domain.Order

// IDGenerator описывает операцию генерации ID.
type IDGenerator interface {
	// NextID возвращает уникальный ID.
	NextID() int64
}

// OrderRepositoryInMemory хранит заказы в in-memory хранилище.
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

// GetByIDOrderItemsBySKU возвращает заказ по идентификатору, если он существует.
func (or *OrderRepositoryInMemory) GetByIDOrderItemsBySKU(_ context.Context, orderID int64) (*domain.Order, error) {
	or.mx.RLock()
	defer or.mx.RUnlock()

	order, ok := or.storage[orderID]
	if !ok {
		return nil, domain.ErrOrderNotExist
	}

	orderCopy := *order
	orderCopy.Items = make([]*domain.OrderItem, 0, len(order.Items))
	for _, item := range order.Items {
		itemCopy := *item
		orderCopy.Items = append(orderCopy.Items, &itemCopy)
	}

	sort.Slice(orderCopy.Items, func(i, j int) bool {
		return orderCopy.Items[i].SkuID < orderCopy.Items[j].SkuID
	})

	return &orderCopy, nil
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
