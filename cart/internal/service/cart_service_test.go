package service

import (
	"context"
	"testing"

	"route256/cart/internal/domain"
	mock "route256/cart/mocks"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

type testComponentCS struct {
	cartRepoMock    *mock.CartRepositoryMock
	productServMock *mock.ProductServiceMock
	lomsServMock    *mock.LomsServiceMock
	cartService     *CartService
}

func newTestComponentCS(t *testing.T) *testComponentCS {
	mc := minimock.NewController(t)
	cartRepoMock := mock.NewCartRepositoryMock(mc)
	prouctServMock := mock.NewProductServiceMock(mc)
	lomsServMock := mock.NewLomsServiceMock(mc)
	cartService := NewCartService(cartRepoMock, prouctServMock, lomsServMock)

	return &testComponentCS{
		cartRepoMock:    cartRepoMock,
		productServMock: prouctServMock,
		cartService:     cartService,
		lomsServMock:    lomsServMock,
	}
}

func TestCartService(t *testing.T) {
	t.Parallel()

	t.Run("add cart item success", func(t *testing.T) {
		t.Parallel()

		tc := newTestComponentCS(t)

		ctx := context.Background()
		item := &domain.CartItem{Sku: 1, Count: 2, Name: "name 1", Price: 100}
		returnedProduct := &domain.Product{Sku: 1, Name: "name 1", Price: 100}
		userID := int64(1)

		tc.productServMock.GetProductBySkuMock.When(ctx, item.Sku).Then(returnedProduct, nil)
		tc.lomsServMock.GetStockInfoMock.When(ctx, item.Sku).Then(100, nil)
		tc.cartRepoMock.UpsertCartItemMock.When(ctx, userID, item).Then(item, nil)

		addedItem, err := tc.cartService.AddCartItem(ctx, userID, item)
		require.NoError(t, err)

		assert.Equal(t, *item, *addedItem)
	})

	t.Run("add cart item with unexisting SKU at product service with error", func(t *testing.T) {
		t.Parallel()

		tc := newTestComponentCS(t)

		ctx := context.Background()
		item := &domain.CartItem{Sku: 1, Count: 2, Name: "name 1", Price: 100}
		userID := int64(1)

		tc.productServMock.GetProductBySkuMock.When(ctx, item.Sku).Then(nil, domain.ErrProductNotFound)

		addedItem, err := tc.cartService.AddCartItem(ctx, userID, item)
		require.Error(t, err)

		assert.Nil(t, addedItem)
	})

	t.Run("get cart with total price success", func(t *testing.T) {
		t.Parallel()

		tc := newTestComponentCS(t)

		ctx := context.Background()
		item1 := &domain.CartItem{Sku: 1, Count: 2}
		item2 := &domain.CartItem{Sku: 2, Count: 2}
		item3 := &domain.CartItem{Sku: 3, Count: 2}
		product1 := &domain.Product{Name: "name 1", Price: 100}
		product2 := &domain.Product{Name: "name 2", Price: 300}
		product3 := &domain.Product{Name: "name 3", Price: 200}
		userID := int64(1)

		tc.cartRepoMock.GetCartByUserIDOrderBySkuMock.
			When(ctx, userID).
			Then(&domain.Cart{Items: []*domain.CartItem{item1, item2, item3}}, nil)

		tc.productServMock.GetProductBySkuMock.When(minimock.AnyContext, item1.Sku).Then(product1, nil)
		tc.productServMock.GetProductBySkuMock.When(minimock.AnyContext, item2.Sku).Then(product2, nil)
		tc.productServMock.GetProductBySkuMock.When(minimock.AnyContext, item3.Sku).Then(product3, nil)

		cart, err := tc.cartService.GetCart(ctx, userID)
		require.NoError(t, err)

		assert.Len(t, cart.Items, 3)
		assert.EqualValues(t, 2*100+2*300+2*200, cart.TotalPrice)
	})

	t.Run("delete item from cart", func(t *testing.T) {
		t.Parallel()

		tc := newTestComponentCS(t)

		ctx := context.Background()
		item := &domain.CartItem{Sku: 1, Count: 2, Name: "name 1", Price: 100}
		userID := int64(1)

		tc.cartRepoMock.DeleteCartItemMock.
			When(ctx, userID, item.Sku).
			Then(nil)

		err := tc.cartService.DeleteCartItem(ctx, userID, item.Sku)
		require.NoError(t, err)
	})

	t.Run("add item to cart with out of stock", func(t *testing.T) {
		t.Parallel()

		tc := newTestComponentCS(t)

		ctx := context.Background()
		item := &domain.CartItem{Sku: 1, Count: 2, Name: "name 1", Price: 100}
		product := &domain.Product{Sku: 1, Name: "name 1", Price: 100}
		userID := int64(1)

		tc.productServMock.GetProductBySkuMock.When(ctx, item.Sku).Then(product, nil)
		tc.lomsServMock.GetStockInfoMock.When(ctx, item.Sku).Then(1, nil)

		_, err := tc.cartService.AddCartItem(ctx, userID, item)
		require.Error(t, err)
	})
}
