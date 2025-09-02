package service

import (
	"context"
	"testing"

	"route256/cart/internal/domain"
	mock "route256/cart/mocks"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testComponentCS struct {
	cartRepoMock    *mock.CartRepositoryMock
	productServMock *mock.ProductServiceMock
	cartService     *CartService
}

func newTestComponent(t *testing.T) *testComponentCS {
	mc := minimock.NewController(t)
	cartRepoMock := mock.NewCartRepositoryMock(mc)
	prouctServMock := mock.NewProductServiceMock(mc)
	cartService := NewCartService(cartRepoMock, prouctServMock)

	return &testComponentCS{
		cartRepoMock:    cartRepoMock,
		productServMock: prouctServMock,
		cartService:     cartService,
	}
}

func TestCartService(t *testing.T) {
	t.Parallel()

	t.Run("add cart item success", func(t *testing.T) {
		t.Parallel()

		tc := newTestComponent(t)

		ctx := context.Background()
		item := &domain.CartItem{Sku: 1, Count: 2, Name: "name 1", Price: 100}
		returnedProduct := &domain.Product{Sku: 1, Name: "name 1", Price: 100}
		userID := int64(1)

		tc.productServMock.GetProductBySkuMock.When(ctx, item.Sku).Then(returnedProduct, nil)
		tc.cartRepoMock.UpsertCartItemMock.When(ctx, userID, item).Then(item, nil)

		addedItem, err := tc.cartService.AddCartItem(ctx, userID, item)
		require.NoError(t, err)

		assert.Equal(t, *item, *addedItem)
	})

	t.Run("add cart item with unexisting SKU at product service with error", func(t *testing.T) {
		t.Parallel()

		tc := newTestComponent(t)

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

		tc := newTestComponent(t)

		ctx := context.Background()
		item1 := &domain.CartItem{Sku: 1, Count: 2, Name: "name 1", Price: 100}
		item2 := &domain.CartItem{Sku: 2, Count: 2, Name: "name 2", Price: 300}
		item3 := &domain.CartItem{Sku: 3, Count: 2, Name: "name 3", Price: 200}
		userID := int64(1)

		tc.cartRepoMock.GetCartByUserIDOrderBySkuMock.
			When(ctx, userID).
			Then(&domain.Cart{Items: []*domain.CartItem{item1, item2, item3}}, nil)

		cart, err := tc.cartService.GetCart(ctx, userID)
		require.NoError(t, err)

		assert.Len(t, cart.Items, 3)
		assert.EqualValues(t, 2*100+2*300+2*200, cart.TotalPrice)
	})
}
