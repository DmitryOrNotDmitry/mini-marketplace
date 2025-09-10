package handler

import (
	"context"
	"testing"

	"route256/loms/internal/domain"
	"route256/loms/pkg/api/orders/v1"

	mock "route256/loms/mocks"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testComponentOH struct {
	orderServMock *mock.OrderServiceMock
	orderHandler  *OrderServerGRPC
}

func newTestComponentOH(t *testing.T) *testComponentOH {
	mc := minimock.NewController(t)
	orderServMock := mock.NewOrderServiceMock(mc)
	orderHandler := NewOrderServerGRPC(orderServMock)

	return &testComponentOH{
		orderServMock: orderServMock,
		orderHandler:  orderHandler,
	}
}

func TestOrderServerGRPC_OrderCreate(t *testing.T) {
	t.Parallel()

	t.Run("create order success", func(t *testing.T) {
		t.Parallel()
		tc := newTestComponentOH(t)

		req := &orders.OrderCreateRequest{
			UserId: 42,
			Items: []*orders.ItemInfo{
				{SkuId: 1001, Count: 2},
			},
		}

		expectedOrder := &domain.Order{
			UserID: 42,
			Items: []*domain.OrderItem{
				{SkuID: 1001, Count: 2},
			},
		}

		orderID := int64(777)

		tc.orderServMock.CreateMock.When(context.Background(), expectedOrder).
			Then(orderID, nil)

		res, err := tc.orderHandler.OrderCreate(context.Background(), req)
		require.NoError(t, err)
		assert.Equal(t, orderID, res.OrderId)
	})

	t.Run("create order failed: can not reserve item", func(t *testing.T) {
		t.Parallel()
		tc := newTestComponentOH(t)

		req := &orders.OrderCreateRequest{
			UserId: 10,
			Items: []*orders.ItemInfo{
				{SkuId: 999, Count: 5},
			},
		}

		tc.orderServMock.CreateMock.Return(0, domain.ErrCanNotReserveItem)

		res, err := tc.orderHandler.OrderCreate(context.Background(), req)
		require.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestOrderServerGRPC_OrderInfo(t *testing.T) {
	t.Parallel()

	t.Run("order info success", func(t *testing.T) {
		t.Parallel()
		tc := newTestComponentOH(t)

		req := &orders.OrderInfoRequest{OrderId: 123}
		expectedOrder := &domain.Order{
			UserID: 50,
			Status: domain.AwaitingPayment,
			Items: []*domain.OrderItem{
				{SkuID: 1001, Count: 3},
			},
		}

		tc.orderServMock.GetInfoByIDMock.When(context.Background(), int64(123)).
			Then(expectedOrder, nil)

		res, err := tc.orderHandler.OrderInfo(context.Background(), req)
		require.NoError(t, err)
		require.NotNil(t, res)
		assert.Equal(t, int64(50), res.UserId)
		assert.EqualValues(t, domain.AwaitingPayment, res.Status)
		require.Len(t, res.Items, 1)
		assert.Equal(t, int64(1001), res.Items[0].SkuId)
	})

	t.Run("order not found", func(t *testing.T) {
		t.Parallel()
		tc := newTestComponentOH(t)

		req := &orders.OrderInfoRequest{OrderId: 404}
		tc.orderServMock.GetInfoByIDMock.Return(nil, domain.ErrOrderNotExist)

		res, err := tc.orderHandler.OrderInfo(context.Background(), req)
		require.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestOrderServerGRPC_OrderPay(t *testing.T) {
	t.Parallel()

	t.Run("pay order success", func(t *testing.T) {
		t.Parallel()
		tc := newTestComponentOH(t)

		req := &orders.OrderPayRequest{OrderId: 101}

		tc.orderServMock.PayByIDMock.Expect(context.Background(), int64(101)).
			Return(nil)

		res, err := tc.orderHandler.OrderPay(context.Background(), req)
		require.NoError(t, err)
		assert.NotNil(t, res)
	})

	t.Run("pay order not found", func(t *testing.T) {
		t.Parallel()
		tc := newTestComponentOH(t)

		req := &orders.OrderPayRequest{OrderId: 404}

		tc.orderServMock.PayByIDMock.Return(domain.ErrOrderNotExist)

		res, err := tc.orderHandler.OrderPay(context.Background(), req)
		require.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("pay order invalid status", func(t *testing.T) {
		t.Parallel()
		tc := newTestComponentOH(t)

		req := &orders.OrderPayRequest{OrderId: 202}

		tc.orderServMock.PayByIDMock.Return(domain.ErrPayWithInvalidOrderStatus)

		res, err := tc.orderHandler.OrderPay(context.Background(), req)
		require.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestOrderServerGRPC_OrderCancel(t *testing.T) {
	t.Parallel()

	t.Run("cancel order success", func(t *testing.T) {
		t.Parallel()
		tc := newTestComponentOH(t)

		req := &orders.OrderCancelRequest{OrderId: 501}

		tc.orderServMock.CancelByIDMock.When(context.Background(), int64(501)).
			Then(nil)

		res, err := tc.orderHandler.OrderCancel(context.Background(), req)
		require.NoError(t, err)
		assert.NotNil(t, res)
	})

	t.Run("cancel order not found", func(t *testing.T) {
		t.Parallel()
		tc := newTestComponentOH(t)

		req := &orders.OrderCancelRequest{OrderId: 404}

		tc.orderServMock.CancelByIDMock.Return(domain.ErrOrderNotExist)

		res, err := tc.orderHandler.OrderCancel(context.Background(), req)
		require.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("cancel order invalid status", func(t *testing.T) {
		t.Parallel()
		tc := newTestComponentOH(t)

		req := &orders.OrderCancelRequest{OrderId: 600}

		tc.orderServMock.CancelByIDMock.Return(domain.ErrCancelWithInvalidOrderStatus)

		res, err := tc.orderHandler.OrderCancel(context.Background(), req)
		require.Error(t, err)
		assert.Nil(t, res)
	})
}
