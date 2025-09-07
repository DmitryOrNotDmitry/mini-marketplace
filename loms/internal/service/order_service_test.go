package service

import (
	"context"
	"route256/loms/internal/domain"
	mock "route256/loms/mocks"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testComponentOS struct {
	orderRepoMock *mock.OrderRepositoryMock
	stockServMock *mock.StockServiceIMock
	orderService  *OrderService
}

func newTestComponentOS(t *testing.T) *testComponentOS {
	mc := minimock.NewController(t)
	orderRepoMock := mock.NewOrderRepositoryMock(mc)
	stockServMock := mock.NewStockServiceIMock(mc)
	orderService := NewOrderService(orderRepoMock, stockServMock)

	return &testComponentOS{
		orderRepoMock: orderRepoMock,
		stockServMock: stockServMock,
		orderService:  orderService,
	}
}

func TestOrderService(t *testing.T) {
	t.Parallel()

	t.Run("create order success", func(t *testing.T) {
		t.Parallel()

		tc := newTestComponentOS(t)

		ctx := context.Background()
		order := &domain.Order{UserID: 1, Items: []*domain.OrderItem{}}
		orderSaved := &domain.Order{UserID: 1, Items: []*domain.OrderItem{}, Status: domain.AwaitingPayment}

		tc.stockServMock.ReserveForMock.Return(nil)
		tc.orderRepoMock.InsertMock.When(ctx, orderSaved).Then(1, nil)

		orderID, err := tc.orderService.Create(ctx, order)
		require.NoError(t, err)

		assert.EqualValues(t, 1, orderID)
	})

	t.Run("create order success with out of stock", func(t *testing.T) {
		t.Parallel()

		tc := newTestComponentOS(t)

		ctx := context.Background()
		order := &domain.Order{UserID: 1, Items: []*domain.OrderItem{}}
		orderSaved := &domain.Order{UserID: 1, Items: []*domain.OrderItem{}, Status: domain.Failed}

		tc.stockServMock.ReserveForMock.Return(domain.ErrCanNotReserveItem)
		tc.orderRepoMock.InsertMock.When(ctx, orderSaved).Then(1, nil)

		_, err := tc.orderService.Create(ctx, order)
		require.Error(t, err)
	})

	t.Run("get order info success", func(t *testing.T) {
		t.Parallel()

		tc := newTestComponentOS(t)

		ctx := context.Background()
		orderID := int64(1)
		orderOut := &domain.Order{OrderID: orderID, UserID: 1, Items: []*domain.OrderItem{}, Status: domain.AwaitingPayment}

		tc.orderRepoMock.GetByIDMock.When(ctx, orderID).Then(orderOut, nil)

		order, err := tc.orderService.GetInfoByID(ctx, orderID)
		require.NoError(t, err)

		assert.Equal(t, orderOut, order)
	})

	t.Run("pay order success", func(t *testing.T) {
		t.Parallel()

		tc := newTestComponentOS(t)

		ctx := context.Background()
		orderID := int64(1)
		orderOut := &domain.Order{OrderID: orderID, UserID: 1, Items: []*domain.OrderItem{}, Status: domain.AwaitingPayment}

		tc.orderRepoMock.GetByIDMock.When(ctx, orderID).Then(orderOut, nil)
		tc.stockServMock.ConfirmReserveForMock.Return(nil)
		tc.orderRepoMock.UpdateStatusMock.When(ctx, orderID, domain.Payed).Then(nil)

		err := tc.orderService.PayByID(ctx, orderID)
		require.NoError(t, err)
	})

	t.Run("pay unexisted order", func(t *testing.T) {
		t.Parallel()

		tc := newTestComponentOS(t)

		ctx := context.Background()
		orderID := int64(1)

		tc.orderRepoMock.GetByIDMock.When(ctx, orderID).Then(nil, domain.ErrOrderNotExist)

		err := tc.orderService.PayByID(ctx, orderID)
		require.Error(t, err)
	})

	t.Run("pay order with wrong statuses", func(t *testing.T) {
		t.Parallel()

		tc := newTestComponentOS(t)

		ctx := context.Background()

		for orderID, status := range []domain.Status{domain.Failed, domain.Cancelled, domain.New} {
			orderOut := &domain.Order{OrderID: int64(orderID), UserID: 1, Items: []*domain.OrderItem{}, Status: status}

			tc.orderRepoMock.GetByIDMock.When(ctx, int64(orderID)).Then(orderOut, nil)

			err := tc.orderService.PayByID(ctx, int64(orderID))
			require.Error(t, err)
		}
	})

	t.Run("pay already payed order", func(t *testing.T) {
		t.Parallel()

		tc := newTestComponentOS(t)

		ctx := context.Background()
		orderID := int64(1)

		orderOut := &domain.Order{OrderID: orderID, UserID: 1, Items: []*domain.OrderItem{}, Status: domain.Payed}

		tc.orderRepoMock.GetByIDMock.When(ctx, orderID).Then(orderOut, nil)

		err := tc.orderService.PayByID(ctx, orderID)
		require.NoError(t, err)
	})

	t.Run("cancel order success", func(t *testing.T) {
		t.Parallel()

		tc := newTestComponentOS(t)

		ctx := context.Background()
		orderID := int64(1)
		orderOut := &domain.Order{OrderID: orderID, UserID: 1, Items: []*domain.OrderItem{}, Status: domain.AwaitingPayment}

		tc.orderRepoMock.GetByIDMock.When(ctx, orderID).Then(orderOut, nil)
		tc.stockServMock.CancelReserveForMock.Return(nil)
		tc.orderRepoMock.UpdateStatusMock.When(ctx, orderID, domain.Cancelled).Then(nil)

		err := tc.orderService.CancelByID(ctx, orderID)
		require.NoError(t, err)
	})

	t.Run("cancel order with wrong statuses", func(t *testing.T) {
		t.Parallel()

		tc := newTestComponentOS(t)

		ctx := context.Background()

		for orderID, status := range []domain.Status{domain.Failed, domain.Payed} {
			orderOut := &domain.Order{OrderID: int64(orderID), UserID: 1, Items: []*domain.OrderItem{}, Status: status}

			tc.orderRepoMock.GetByIDMock.When(ctx, int64(orderID)).Then(orderOut, nil)

			err := tc.orderService.CancelByID(ctx, int64(orderID))
			require.Error(t, err)
		}
	})

	t.Run("cancel already canceled order", func(t *testing.T) {
		t.Parallel()

		tc := newTestComponentOS(t)

		ctx := context.Background()
		orderID := int64(1)

		orderOut := &domain.Order{OrderID: orderID, UserID: 1, Items: []*domain.OrderItem{}, Status: domain.Cancelled}

		tc.orderRepoMock.GetByIDMock.When(ctx, orderID).Then(orderOut, nil)

		err := tc.orderService.CancelByID(ctx, orderID)
		require.NoError(t, err)
	})

}
