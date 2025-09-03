package service

import (
	"context"
	"net/http"
	"net/http/httptest"
	"route256/cart/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProductServiceHttp_GetProductBySku(t *testing.T) {
	t.Parallel()

	t.Run("successful response", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(`{"name":"TestProduct","price":100,"sku":12345}`))
			require.NoError(t, err)
		}))
		defer server.Close()

		svc := NewProductServiceHTTP(http.Client{}, "token", server.URL)

		got, err := svc.GetProductBySku(context.Background(), 12345)
		require.NoError(t, err)

		expected := &domain.Product{Name: "TestProduct", Price: 100, Sku: 12345}
		assert.Equal(t, expected, got)
	})

	t.Run("product not found", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			_, err := w.Write([]byte(``))
			require.NoError(t, err)
		}))
		defer server.Close()

		svc := NewProductServiceHTTP(http.Client{}, "token", server.URL)

		got, err := svc.GetProductBySku(context.Background(), 99999)
		require.Error(t, err)

		assert.Nil(t, got)
	})
}
