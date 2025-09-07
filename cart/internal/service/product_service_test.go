package service

import (
	"context"
	"io"
	"net/http"
	"route256/cart/internal/domain"
	mock "route256/cart/mocks"
	"strings"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testComponentPS struct {
	httpClientMock *mock.HTTPClientMock
	productService *ProductServiceHTTP
}

func newTestComponentPS(t *testing.T) *testComponentPS {
	mc := minimock.NewController(t)
	httpClientMock := mock.NewHTTPClientMock(mc)
	productService := NewProductServiceHTTP(httpClientMock, "token", "url-test")

	return &testComponentPS{
		httpClientMock: httpClientMock,
		productService: productService,
	}
}

func TestProductServiceHttp_GetProductBySku(t *testing.T) {
	t.Parallel()

	t.Run("successful response", func(t *testing.T) {
		t.Parallel()

		tc := newTestComponentPS(t)

		body := io.NopCloser(strings.NewReader(`{"name":"TestProduct","price":100,"sku":12345}`))
		response := &http.Response{
			StatusCode: http.StatusOK,
			Body:       body,
		}

		tc.httpClientMock.DoMock.Return(response, nil)

		product, err := tc.productService.GetProductBySku(context.Background(), 12345)
		require.NoError(t, err)

		expected := &domain.Product{Name: "TestProduct", Price: 100, Sku: 12345}
		assert.Equal(t, expected, product)
	})

	t.Run("product not found", func(t *testing.T) {
		t.Parallel()

		tc := newTestComponentPS(t)

		response := &http.Response{
			StatusCode: http.StatusNotFound,
			Body:       http.NoBody,
		}

		tc.httpClientMock.DoMock.Return(response, nil)

		product, err := tc.productService.GetProductBySku(context.Background(), 99999)
		require.Error(t, err)

		assert.Nil(t, product)
	})
}
