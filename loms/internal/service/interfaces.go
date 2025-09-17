package service

import (
	"context"
	"route256/loms/internal/domain"
)

// TxManager определяет интерфейс для работы с транзакциями.
type TxManager interface {
	// WithTransaction выполняет функцию fn в дефолтном уровне изоляции транзакции.
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) (err error)
}

// OrderRepository описывает методы работы с заказами в хранилище.
type OrderRepository interface {
	// Insert добавляет новый заказ и возвращает его ID.
	Insert(ctx context.Context, order *domain.Order) (int64, error)
	// GetByIDOrderItemsBySKU возвращает заказ с деталями по его ID и сортирует товары по SKU.
	GetByIDOrderItemsBySKU(ctx context.Context, orderID int64) (*domain.Order, error)
	// UpdateStatus обновляет статус заказа.
	UpdateStatus(ctx context.Context, orderID int64, newStatus domain.Status) error
}

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
