package service

import (
	"context"
	"route256/loms/internal/domain"
)

// OperationType определяет тип операции.
type OperationType int

const (
	// Read операция только для чтения.
	Read OperationType = 1
	// Write операция на запись.
	Write OperationType = 2
	// FromTx операция извлекается из транзакции.
	FromTx OperationType = 3
)

// TxManager определяет интерфейс для работы с транзакциями.
type TxManager interface {
	// WithTransaction выполняет функцию fn в дефолтном уровне изоляции транзакции.
	WithTransaction(ctx context.Context, operationType OperationType, fn func(ctx context.Context) error) (err error)

	// WithRepeatableRead выполняет функцию fn в RepeatableRead уровне изоляции транзакции.
	WithRepeatableRead(ctx context.Context, operationType OperationType, fn func(ctx context.Context) error) (err error)
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
	// GetBySkuIDForUpdate возвращает информацию о запасе по SKU с блокировкой на обновление.
	GetBySkuIDForUpdate(ctx context.Context, skuID int64) (*domain.Stock, error)
}

// OrderEventRepository описывает методы работы с событиями о заказе.
type OrderEventRepository interface {
	// Insert добавляет новое событие по статусу в заказе
	Insert(ctx context.Context, order *domain.Order) error

	// GetUnprocessedEventsLimit возвращает неотправленные события с ограничением по количеству.
	GetUnprocessedEventsLimit(ctx context.Context, limit int) ([]*domain.OrderEventOutbox, error)
	// UpdateEventStatus обновляет статус события заказа.
	UpdateEventStatus(ctx context.Context, eventID int64, newStatus domain.EventStatus) error
}
