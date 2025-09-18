package postgres

import (
	"context"
	"errors"
	"fmt"
	"route256/cart/pkg/logger"
	"route256/loms/internal/domain"
	repo_sqlc "route256/loms/internal/infra/repository/postgres/sqlc/generated"

	"github.com/jackc/pgx/v5"
)

// NewOrderRepository создает новый OrderRepository.
func NewOrderRepository(pool repo_sqlc.DBTX) *OrderRepository {
	return &OrderRepository{
		repo_sqlc.New(pool),
	}
}

// OrderRepository предоставляет доступ к хранилищу заказов из postgres.
type OrderRepository struct {
	querier repo_sqlc.Querier
}

// Insert добавляет новый заказ и возвращает его ID из postgres.
func (or *OrderRepository) Insert(ctx context.Context, order *domain.Order) (int64, error) {
	orderID, err := or.querier.AddOrder(ctx, &repo_sqlc.AddOrderParams{
		UserID: order.UserID,
		Status: string(order.Status),
	})
	if err != nil {
		return 0, fmt.Errorf("querier.AddOrder: %w", err)
	}

	for _, item := range order.Items {
		err = or.querier.AddOrderItem(ctx, &repo_sqlc.AddOrderItemParams{
			Sku:     item.SkuID,
			OrderID: orderID,
			Count:   int64(item.Count),
		})
		if err != nil {
			return 0, fmt.Errorf("querier.AddOrderItem: %w", err)
		}
	}

	return orderID, nil
}

// GetByIDOrderItemsBySKU возвращает заказ с деталями по его ID и сортирует товары по SKU из postgres.
func (or *OrderRepository) GetByIDOrderItemsBySKU(ctx context.Context, orderID int64) (*domain.Order, error) {
	orderDB, err := or.querier.GetOrderByID(ctx, orderID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrOrderNotExist
		}

		return nil, fmt.Errorf("querier.GetOrderByID: %w", err)
	}

	orderItemsDB, err := or.querier.GetOrderItemsOrderBySKU(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("querier.GetOrderItemsOrderBySKU: %w", err)
	}

	order := &domain.Order{
		OrderID: orderDB.OrderID,
		UserID:  orderDB.UserID,
		Status:  domain.Status(orderDB.Status),
		Items:   make([]*domain.OrderItem, 0, len(orderItemsDB)),
	}
	for _, itemDB := range orderItemsDB {
		count, err := Int64ToUint32(itemDB.Count)
		if err != nil {
			logger.Warning(fmt.Sprintf("Int64ToUint32 (Count=%d): %s", itemDB.Count, err.Error()))
		}

		order.Items = append(order.Items, &domain.OrderItem{
			SkuID: itemDB.Sku,
			Count: count,
		})
	}

	return order, nil
}

// UpdateStatus обновляет статус заказа из postgres.
func (or *OrderRepository) UpdateStatus(ctx context.Context, orderID int64, newStatus domain.Status) error {
	err := or.querier.UpdateStatusByID(ctx, &repo_sqlc.UpdateStatusByIDParams{
		OrderID: orderID,
		Status:  string(newStatus),
	})
	if err != nil {
		return fmt.Errorf("querier.UpdateStatusByID: %w", err)
	}

	return nil
}
