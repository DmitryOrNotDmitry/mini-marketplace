package handler

import (
	"context"
	"errors"
	"route256/loms/internal/domain"
	"route256/loms/pkg/api/orders/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrderService interface {
	Create(ctx context.Context, order *domain.Order) (int64, error)
	GetInfoByID(ctx context.Context, orderID int64) (*domain.Order, error)
	PayByID(ctx context.Context, orderID int64) error
	CancelByID(ctx context.Context, orderID int64) error
}

type OrderServerGRPC struct {
	orderService OrderService
	orders.UnimplementedOrderServiceServer
}

// NewOrderServerGRPC создает новый экземпляр OrderServerGRPC.
func NewOrderServerGRPC(orderService OrderService) *OrderServerGRPC {
	return &OrderServerGRPC{
		orderService: orderService,
	}
}

// OrderCreate создает новый заказ.
func (os *OrderServerGRPC) OrderCreate(ctx context.Context, req *orders.OrderCreateRequest) (*orders.OrderCreateResponse, error) {
	order := &domain.Order{
		UserID: req.UserId,
		Items:  make([]*domain.OrderItem, 0, len(req.Items)),
	}
	for _, reqItem := range req.Items {
		order.Items = append(order.Items, &domain.OrderItem{
			SkuID: reqItem.SkuId,
			Count: reqItem.Count,
		})
	}

	orderID, err := os.orderService.Create(ctx, order)
	if err != nil {
		if errors.Is(err, domain.ErrItemStockNotExist) ||
			errors.Is(err, domain.ErrCanNotReserveItem) {
			return nil, status.Errorf(codes.FailedPrecondition, err.Error())
		}

		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	res := &orders.OrderCreateResponse{OrderId: orderID}

	return res, nil
}

// OrderInfo возвращает информацию о заказе по его идентификатору.
func (os *OrderServerGRPC) OrderInfo(ctx context.Context, req *orders.OrderInfoRequest) (*orders.OrderInfoResponse, error) {
	order, err := os.orderService.GetInfoByID(ctx, req.OrderId)
	if err != nil {
		if errors.Is(err, domain.ErrOrderNotExist) {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}

		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	res := &orders.OrderInfoResponse{
		UserId: order.UserID,
		Status: string(order.Status),
		Items:  make([]*orders.ItemInfo, 0, len(order.Items)),
	}
	for _, item := range order.Items {
		res.Items = append(res.Items, &orders.ItemInfo{
			SkuId: item.SkuID,
			Count: item.Count,
		})
	}

	return res, nil
}

// OrderPay помечает заказ как оплаченный по его идентификатору.
func (os *OrderServerGRPC) OrderPay(ctx context.Context, req *orders.OrderPayRequest) (*orders.OrderPayResponse, error) {
	err := os.orderService.PayByID(ctx, req.OrderId)
	if err != nil {
		if errors.Is(err, domain.ErrOrderNotExist) {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}

		if errors.Is(err, domain.ErrPayWithInvalidOrderStatus) {
			return nil, status.Errorf(codes.FailedPrecondition, err.Error())
		}

		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	return &orders.OrderPayResponse{}, nil
}

// OrderCancel отменяет заказ по его идентификатору.
func (os *OrderServerGRPC) OrderCancel(ctx context.Context, req *orders.OrderCancelRequest) (*orders.OrderCancelResponse, error) {
	err := os.orderService.CancelByID(ctx, req.OrderId)
	if err != nil {
		if errors.Is(err, domain.ErrOrderNotExist) {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}

		if errors.Is(err, domain.ErrCancelWithInvalidOrderStatus) {
			return nil, status.Errorf(codes.FailedPrecondition, err.Error())
		}

		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	return &orders.OrderCancelResponse{}, nil
}
