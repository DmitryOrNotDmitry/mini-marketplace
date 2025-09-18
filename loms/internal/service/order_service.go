package service

import (
	"context"
	"fmt"
	"route256/loms/internal/domain"
)

// OrderRepoFactory создает экземпляры репозитория заказов.
type OrderRepoFactory interface {
	// CreateOrder создает новый репозиторий заказов.
	CreateOrder(ctx context.Context, operationType OperationType) OrderRepository
}

// StockServiceI описывает методы работы с резервированием товаров.
type StockServiceI interface {
	// ReserveFor резервирует товары под заказ.
	ReserveFor(ctx context.Context, order *domain.Order) error
	// CancelReserveFor отменяет резервирование товаров по заказу.
	CancelReserveFor(ctx context.Context, order *domain.Order) error
	// ConfirmReserveFor подтверждает резервирование и уменьшает общий запас.
	ConfirmReserveFor(ctx context.Context, order *domain.Order) error
}

// OrderService реализует бизнес-логику управления заказами.
type OrderService struct {
	stockService      StockServiceI
	repositoryFactory OrderRepoFactory
	txManager         TxManager
}

// NewOrderService создает новый сервис управления заказами.
func NewOrderService(stockService StockServiceI, repositoryFactory OrderRepoFactory, txManager TxManager) *OrderService {
	return &OrderService{
		stockService:      stockService,
		repositoryFactory: repositoryFactory,
		txManager:         txManager,
	}
}

// Create создает новый заказ, резервирует товары и возвращает идентификатор заказа.
func (os *OrderService) Create(ctx context.Context, order *domain.Order) (int64, error) {
	var orderID int64
	var stockErr error

	err := os.txManager.WithTransaction(ctx, Write, func(ctx context.Context) error {
		order.Status = domain.New

		stockErr = os.stockService.ReserveFor(ctx, order)
		if stockErr != nil {
			stockErr = fmt.Errorf("stockService.ReserveFor: %w", stockErr)
			order.Status = domain.Failed
		} else {
			order.Status = domain.AwaitingPayment
		}

		orderRepository := os.repositoryFactory.CreateOrder(ctx, FromTx)
		var err error
		orderID, err = orderRepository.Insert(ctx, order)
		if err != nil {
			return fmt.Errorf("orderRepository.Insert: %w", err)
		}
		return nil
	})
	if err != nil {
		return 0, fmt.Errorf("txManager.WithTransaction: %w", err)
	}

	return orderID, stockErr
}

// GetInfoByID возвращает информацию о заказе по его идентификатору.
func (os *OrderService) GetInfoByID(ctx context.Context, orderID int64) (*domain.Order, error) {
	orderRepository := os.repositoryFactory.CreateOrder(ctx, Read)
	order, err := orderRepository.GetByIDOrderItemsBySKU(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("orderRepository.GetByIDOrderItemsBySKU: %w", err)
	}

	return order, nil
}

// PayByID подтверждает оплату заказа по идентификатору.
func (os *OrderService) PayByID(ctx context.Context, orderID int64) error {
	err := os.txManager.WithRepeatableRead(ctx, Write, func(ctx context.Context) error {
		order, err := os.GetInfoByID(ctx, orderID)
		if err != nil {
			return fmt.Errorf("os.GetInfoByID: %w", err)
		}

		if order.Status == domain.Paid {
			return nil
		}

		if order.Status != domain.AwaitingPayment {
			return domain.ErrPayWithInvalidOrderStatus
		}

		err = os.stockService.ConfirmReserveFor(ctx, order)
		if err != nil {
			return fmt.Errorf("stockService.ConfirmReserveFor: %w", err)
		}

		orderRepository := os.repositoryFactory.CreateOrder(ctx, FromTx)
		return orderRepository.UpdateStatus(ctx, orderID, domain.Paid)
	})
	if err != nil {
		return fmt.Errorf("txManager.WithRepeatableRead: %w", err)
	}

	return nil
}

// CancelByID отменяет заказ по идентификатору.
func (os *OrderService) CancelByID(ctx context.Context, orderID int64) error {
	err := os.txManager.WithRepeatableRead(ctx, Write, func(ctx context.Context) error {
		order, err := os.GetInfoByID(ctx, orderID)
		if err != nil {
			return fmt.Errorf("os.GetInfoByID: %w", err)
		}

		if order.Status == domain.Cancelled {
			return nil
		}

		if order.Status == domain.Paid || order.Status == domain.Failed {
			return domain.ErrCancelWithInvalidOrderStatus
		}

		errCancel := os.stockService.CancelReserveFor(ctx, order)
		if errCancel != nil {
			return fmt.Errorf("stockService.CancelReserveFor: %s", errCancel)
		}

		orderRepository := os.repositoryFactory.CreateOrder(ctx, FromTx)
		return orderRepository.UpdateStatus(ctx, orderID, domain.Cancelled)
	})
	if err != nil {
		return fmt.Errorf("txManager.WithRepeatableRead: %w", err)
	}

	return nil
}
