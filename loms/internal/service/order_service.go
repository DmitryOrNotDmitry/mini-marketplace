package service

import (
	"context"
	"fmt"
	"route256/cart/pkg/logger"
	"route256/loms/internal/domain"
)

type OrderRepository interface {
	Insert(ctx context.Context, order *domain.Order) (int64, error)
	GetByID(ctx context.Context, orderID int64) (*domain.Order, error)
	UpdateStatus(ctx context.Context, orderID int64, newStatus domain.Status) error
}

type StockServiceI interface {
	ReserveFor(ctx context.Context, order *domain.Order) error
	CancelReserveFor(ctx context.Context, order *domain.Order) error
	ConfirmReserveFor(ctx context.Context, order *domain.Order) error
}

type OrderService struct {
	orderRepository OrderRepository
	stockService    StockServiceI
}

// NewOrderService создает новый сервис управления заказами.
func NewOrderService(orderRepository OrderRepository, stockService StockServiceI) *OrderService {
	return &OrderService{
		orderRepository: orderRepository,
		stockService:    stockService,
	}
}

// Create создает новый заказ, резервирует товары и возвращает идентификатор заказа.
func (os *OrderService) Create(ctx context.Context, order *domain.Order) (int64, error) {
	order.Status = domain.New

	err := os.stockService.ReserveFor(ctx, order)
	if err != nil {
		order.Status = domain.Failed
	} else {
		order.Status = domain.AwaitingPayment
	}

	orderID, errInsert := os.orderRepository.Insert(ctx, order)
	if errInsert != nil {
		return -1, fmt.Errorf("orderRepository.Insert: %w", errInsert)
	}

	return orderID, err
}

// GetInfoByID возвращает информацию о заказе по его идентификатору.
func (os *OrderService) GetInfoByID(ctx context.Context, orderID int64) (*domain.Order, error) {
	order, err := os.orderRepository.GetByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("orderRepository.Insert: %w", err)
	}

	return order, nil
}

// PayByID подтверждает оплату заказа по идентификатору.
func (os *OrderService) PayByID(ctx context.Context, orderID int64) error {
	order, err := os.GetInfoByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("os.GetInfoByID: %w", err)
	}

	if order.Status == domain.Payed {
		return nil
	}

	if order.Status != domain.AwaitingPayment {
		return domain.ErrPayWithInvalidOrderStatus
	}

	errConfirm := os.stockService.ConfirmReserveFor(ctx, order)
	if errConfirm != nil {
		logger.Error(fmt.Sprintf("stockService.ConfirmReserveFor: %s", errConfirm.Error()))
	}

	return os.orderRepository.UpdateStatus(ctx, orderID, domain.Payed)
}

// CancelByID отменяет заказ по идентификатору.
func (os *OrderService) CancelByID(ctx context.Context, orderID int64) error {
	order, err := os.GetInfoByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("os.GetInfoByID: %w", err)
	}

	if order.Status == domain.Cancelled {
		return nil
	}

	if order.Status == domain.Payed || order.Status == domain.Failed {
		return domain.ErrCancelWithInvalidOrderStatus
	}

	errCancel := os.stockService.CancelReserveFor(ctx, order)
	if errCancel != nil {
		logger.Error(fmt.Sprintf("stockService.CancelReserveFor: %s", errCancel.Error()))
	}

	return os.orderRepository.UpdateStatus(ctx, orderID, domain.Cancelled)
}
