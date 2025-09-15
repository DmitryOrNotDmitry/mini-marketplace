package handler

import (
	"context"
	"route256/loms/internal/domain"
	"route256/loms/mocks"
	"route256/loms/pkg/api/stocks/v1"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testComponentSS struct {
	stockServMock *mocks.StockServiceMock
	stockHandler  *StockServerGRPC
}

func newTestComponentSS(t *testing.T) *testComponentSS {
	mc := minimock.NewController(t)
	stockServMock := mocks.NewStockServiceMock(mc)
	stockHandler := NewStockServerGRPC(stockServMock)

	return &testComponentSS{
		stockServMock: stockServMock,
		stockHandler:  stockHandler,
	}
}

func TestStockServerGRPC_StockInfo(t *testing.T) {
	t.Parallel()

	t.Run("stock info success", func(t *testing.T) {
		t.Parallel()
		tc := newTestComponentSS(t)

		req := &stocks.StockInfoRequest{SkuId: 1001}
		tc.stockServMock.GetAvailableCountMock.When(context.Background(), int64(1001)).
			Then(uint32(50), nil)

		res, err := tc.stockHandler.StockInfo(context.Background(), req)
		require.NoError(t, err)
		require.NotNil(t, res)
		assert.Equal(t, uint32(50), res.Count)
	})

	t.Run("stock not exist", func(t *testing.T) {
		t.Parallel()
		tc := newTestComponentSS(t)

		req := &stocks.StockInfoRequest{SkuId: 9999}
		tc.stockServMock.GetAvailableCountMock.Return(uint32(0), domain.ErrItemStockNotExist)

		res, err := tc.stockHandler.StockInfo(context.Background(), req)
		require.Error(t, err)
		assert.Nil(t, res)
	})
}
