package service

import (
	"context"
	"testing"

	"route256/loms/internal/domain"
	mock "route256/loms/mocks"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testComponentCS struct {
	stockRepoMock *mock.StockRepositoryMock
	stockService  *StockService
}

func newTestComponentSS(t *testing.T) *testComponentCS {
	mc := minimock.NewController(t)
	stockRepoMock := mock.NewStockRepositoryMock(mc)
	stockService := NewStockService(stockRepoMock)

	return &testComponentCS{
		stockRepoMock: stockRepoMock,
		stockService:  stockService,
	}
}

func TestStockService(t *testing.T) {
	t.Parallel()

	t.Run("get avaliable count", func(t *testing.T) {
		t.Parallel()

		tc := newTestComponentSS(t)

		ctx := context.Background()
		skuID := int64(1)
		stock := &domain.Stock{SkuID: skuID, TotalCount: 100, Reserved: 10}

		tc.stockRepoMock.GetBySkuIDMock.When(ctx, skuID).Then(stock, nil)

		count, err := tc.stockService.GetAvailableCount(ctx, skuID)
		require.NoError(t, err)

		assert.EqualValues(t, 90, count)
	})

	t.Run("reserve stocks for order", func(t *testing.T) {
		t.Parallel()

		tc := newTestComponentSS(t)

		ctx := context.Background()
		order := &domain.Order{
			OrderID: 1,
			UserID:  1,
			Status:  "New",
			Items: []*domain.OrderItem{
				&domain.OrderItem{SkuID: 1, Count: 100},
				&domain.OrderItem{SkuID: 2, Count: 100},
				&domain.OrderItem{SkuID: 3, Count: 100},
			},
		}

		tc.stockRepoMock.AddReserveMock.Return(nil)

		err := tc.stockService.ReserveFor(ctx, order)
		require.NoError(t, err)
	})

	t.Run("reserve stocks for order with reserve error", func(t *testing.T) {
		t.Parallel()

		tc := newTestComponentSS(t)

		ctx := context.Background()
		order := &domain.Order{
			OrderID: 1,
			UserID:  1,
			Status:  "New",
			Items: []*domain.OrderItem{
				&domain.OrderItem{SkuID: 1, Count: 100},
				&domain.OrderItem{SkuID: 2, Count: 100},
				&domain.OrderItem{SkuID: 3, Count: 100},
			},
		}

		tc.stockRepoMock.AddReserveMock.When(ctx, order.Items[0].SkuID, order.Items[0].Count).Then(nil)
		tc.stockRepoMock.AddReserveMock.When(ctx, order.Items[1].SkuID, order.Items[1].Count).Then(domain.ErrCanNotReserveItem)

		tc.stockRepoMock.RemoveReserveMock.When(ctx, order.Items[0].SkuID, order.Items[0].Count).Then(nil)

		err := tc.stockService.ReserveFor(ctx, order)
		require.Error(t, err)
	})

	t.Run("cancel stocks for order", func(t *testing.T) {
		t.Parallel()

		tc := newTestComponentSS(t)

		ctx := context.Background()
		order := &domain.Order{
			OrderID: 1,
			UserID:  1,
			Status:  "New",
			Items: []*domain.OrderItem{
				&domain.OrderItem{SkuID: 1, Count: 100},
				&domain.OrderItem{SkuID: 2, Count: 100},
				&domain.OrderItem{SkuID: 3, Count: 100},
			},
		}

		tc.stockRepoMock.RemoveReserveMock.Times(3).Return(nil)

		err := tc.stockService.CancelReserveFor(ctx, order)
		require.NoError(t, err)
	})

	t.Run("confirm stocks for order", func(t *testing.T) {
		t.Parallel()

		tc := newTestComponentSS(t)

		ctx := context.Background()
		order := &domain.Order{
			OrderID: 1,
			UserID:  1,
			Status:  "New",
			Items: []*domain.OrderItem{
				&domain.OrderItem{SkuID: 1, Count: 100},
				&domain.OrderItem{SkuID: 2, Count: 100},
				&domain.OrderItem{SkuID: 3, Count: 100},
			},
		}

		tc.stockRepoMock.ReduceReserveAndTotalMock.When(ctx, order.Items[0].SkuID, order.Items[0].Count).Then(nil)
		tc.stockRepoMock.ReduceReserveAndTotalMock.When(ctx, order.Items[1].SkuID, order.Items[1].Count).Then(nil)
		tc.stockRepoMock.ReduceReserveAndTotalMock.When(ctx, order.Items[2].SkuID, order.Items[2].Count).Then(nil)

		err := tc.stockService.ConfirmReserveFor(ctx, order)
		require.NoError(t, err)
	})

}
