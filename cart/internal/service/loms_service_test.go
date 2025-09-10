package service

import (
	"context"
	"errors"
	"testing"

	"route256/cart/internal/domain"
	mock "route256/cart/mocks"
	"route256/loms/pkg/api/orders/v1"
	"route256/loms/pkg/api/stocks/v1"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testComponentLS struct {
	stockClientMock *mock.StockServiceClientMock
	orderClientMock *mock.OrderServiceClientMock
	lomsService     *LomsServiceGRPC
}

func newTestComponentLS(t *testing.T) *testComponentLS {
	mc := minimock.NewController(t)
	stockClientMock := mock.NewStockServiceClientMock(mc)
	orderClientMock := mock.NewOrderServiceClientMock(mc)
	lomsService := NewLomsServiceGRPC(stockClientMock, orderClientMock)

	return &testComponentLS{
		stockClientMock: stockClientMock,
		orderClientMock: orderClientMock,
		lomsService:     lomsService,
	}
}

func TestLomsService(t *testing.T) {
	t.Parallel()

	t.Run("get stock info success", func(t *testing.T) {
		t.Parallel()

		tc := newTestComponentLS(t)

		ctx := context.Background()
		skuID := int64(1)

		tc.stockClientMock.StockInfoMock.When(minimock.AnyContext, &stocks.StockInfoRequest{SkuId: skuID}).
			Then(&stocks.StockInfoResponse{Count: 100}, nil)

		count, err := tc.lomsService.GetStockInfo(ctx, skuID)
		require.NoError(t, err)

		assert.EqualValues(t, 100, count)
	})

	t.Run("get stock info failed: failed stock client", func(t *testing.T) {
		t.Parallel()

		tc := newTestComponentLS(t)

		ctx := context.Background()

		tc.stockClientMock.StockInfoMock.
			Return(nil, errors.New("error"))

		_, err := tc.lomsService.GetStockInfo(ctx, 1)
		require.Error(t, err)
	})

	t.Run("order create success", func(t *testing.T) {
		t.Parallel()

		tc := newTestComponentLS(t)

		ctx := context.Background()
		userID := int64(1)
		cart := &domain.Cart{Items: []*domain.CartItem{
			&domain.CartItem{Sku: 1, Count: 10},
		}}

		tc.orderClientMock.OrderCreateMock.
			Return(&orders.OrderCreateResponse{OrderId: 1}, nil)

		orderID, err := tc.lomsService.OrderCreate(ctx, userID, cart)
		require.NoError(t, err)

		assert.EqualValues(t, 1, orderID)
	})

	t.Run("order create failed: failed order client", func(t *testing.T) {
		t.Parallel()

		tc := newTestComponentLS(t)

		ctx := context.Background()
		userID := int64(1)
		cart := &domain.Cart{Items: []*domain.CartItem{
			&domain.CartItem{Sku: 1, Count: 10},
		}}

		tc.orderClientMock.OrderCreateMock.
			Return(nil, errors.New("error"))

		_, err := tc.lomsService.OrderCreate(ctx, userID, cart)
		require.Error(t, err)
	})
}
