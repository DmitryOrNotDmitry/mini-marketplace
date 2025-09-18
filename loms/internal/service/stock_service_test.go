package service_test

import (
	"context"
	"testing"

	"route256/loms/internal/domain"
	"route256/loms/internal/service"
	mock "route256/loms/mocks"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TxManagerForTests struct {
}

func (tx *TxManagerForTests) WithTransaction(ctx context.Context, operationType service.OperationType, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

func (tx *TxManagerForTests) WithRepeatableRead(ctx context.Context, operationType service.OperationType, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

type testComponentSS struct {
	stockRepoMock   *mock.StockRepositoryMock
	repoFactoryMock *mock.StockRepoFactoryMock
	stockService    *service.StockService
}

func newTestComponentSS(t *testing.T) *testComponentSS {
	mc := minimock.NewController(t)
	stockRepoMock := mock.NewStockRepositoryMock(mc)
	repoFactoryMock := mock.NewStockRepoFactoryMock(mc)
	stockService := service.NewStockService(repoFactoryMock, &TxManagerForTests{})

	return &testComponentSS{
		stockRepoMock:   stockRepoMock,
		stockService:    stockService,
		repoFactoryMock: repoFactoryMock,
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

		tc.repoFactoryMock.CreateStockMock.Return(tc.stockRepoMock)
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
		stock := &domain.Stock{
			Reserved:   uint32(0),
			TotalCount: uint32(100),
		}

		tc.repoFactoryMock.CreateStockMock.Return(tc.stockRepoMock)
		tc.stockRepoMock.GetBySkuIDMock.Return(stock, nil)
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
		stock := &domain.Stock{
			Reserved:   uint32(0),
			TotalCount: uint32(100),
		}

		tc.repoFactoryMock.CreateStockMock.Return(tc.stockRepoMock)
		tc.stockRepoMock.GetBySkuIDMock.Return(stock, nil)
		tc.stockRepoMock.AddReserveMock.When(ctx, order.Items[0].SkuID, order.Items[0].Count).Then(nil)
		tc.stockRepoMock.AddReserveMock.When(ctx, order.Items[1].SkuID, order.Items[1].Count).Then(domain.ErrCanNotReserveItem)

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

		tc.repoFactoryMock.CreateStockMock.Return(tc.stockRepoMock)
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

		tc.repoFactoryMock.CreateStockMock.Return(tc.stockRepoMock)
		tc.stockRepoMock.ReduceReserveAndTotalMock.When(ctx, order.Items[0].SkuID, order.Items[0].Count).Then(nil)
		tc.stockRepoMock.ReduceReserveAndTotalMock.When(ctx, order.Items[1].SkuID, order.Items[1].Count).Then(nil)
		tc.stockRepoMock.ReduceReserveAndTotalMock.When(ctx, order.Items[2].SkuID, order.Items[2].Count).Then(nil)

		err := tc.stockService.ConfirmReserveFor(ctx, order)
		require.NoError(t, err)
	})

}
