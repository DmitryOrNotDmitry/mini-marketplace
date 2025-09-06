package repository

import (
	"context"
	"crypto/rand"
	"math/big"
	"testing"

	"route256/cart/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCartRepositoryInMemory(t *testing.T) {
	t.Parallel()

	t.Run("add new item", func(t *testing.T) {
		t.Parallel()

		repo := NewInMemoryCartRepository(10)
		ctx := context.Background()
		item := &domain.CartItem{
			Sku:   1,
			Name:  "Item 1",
			Count: 1,
			Price: 100,
		}
		userID := int64(1)

		_, err := repo.UpsertCartItem(ctx, userID, item)
		require.NoError(t, err)

		cart, err := repo.GetCartByUserIDOrderBySku(ctx, userID)
		require.NoError(t, err)

		require.Len(t, cart.Items, 1)
		assert.Equal(t, item, cart.Items[0])
	})

	t.Run("add exisiting item", func(t *testing.T) {
		t.Parallel()

		repo := NewInMemoryCartRepository(10)
		ctx := context.Background()
		item := &domain.CartItem{
			Sku:   1,
			Name:  "Item 1",
			Count: 1,
			Price: 100,
		}
		userID := int64(1)

		_, err := repo.UpsertCartItem(ctx, userID, item)
		require.NoError(t, err)
		_, err = repo.UpsertCartItem(ctx, userID, item)
		require.NoError(t, err)

		cart, err := repo.GetCartByUserIDOrderBySku(ctx, userID)
		require.NoError(t, err)

		require.Len(t, cart.Items, 1)
		assert.Equal(t, item.Sku, cart.Items[0].Sku)
		assert.EqualValues(t, 2, cart.Items[0].Count)
	})

	t.Run("delete cartItem", func(t *testing.T) {
		t.Parallel()

		repo := NewInMemoryCartRepository(10)
		ctx := context.Background()
		item := &domain.CartItem{
			Sku:   1,
			Name:  "Item 1",
			Count: 1,
			Price: 100,
		}
		userID := int64(1)

		_, err := repo.UpsertCartItem(ctx, userID, item)
		require.NoError(t, err)

		err = repo.DeleteCartItem(ctx, userID, item.Sku)
		require.NoError(t, err)

		cart, err := repo.GetCartByUserIDOrderBySku(ctx, userID)
		require.NoError(t, err)

		require.Empty(t, cart.Items)
	})

	t.Run("delete other cartItem", func(t *testing.T) {
		t.Parallel()

		repo := NewInMemoryCartRepository(10)
		ctx := context.Background()
		item := &domain.CartItem{
			Sku:   1,
			Name:  "Item 1",
			Count: 1,
			Price: 100,
		}
		userID := int64(1)

		_, err := repo.UpsertCartItem(ctx, userID, item)
		require.NoError(t, err)

		err = repo.DeleteCartItem(ctx, userID, 2)
		require.NoError(t, err)

		cart, err := repo.GetCartByUserIDOrderBySku(ctx, userID)
		require.NoError(t, err)

		require.Len(t, cart.Items, 1)
	})

	t.Run("clear cart", func(t *testing.T) {
		t.Parallel()

		repo := NewInMemoryCartRepository(10)
		ctx := context.Background()
		item := &domain.CartItem{
			Sku:   1,
			Name:  "Item 1",
			Count: 1,
			Price: 100,
		}
		userID := int64(1)

		_, err := repo.UpsertCartItem(ctx, userID, item)
		require.NoError(t, err)

		err = repo.DeleteCart(ctx, userID)
		require.NoError(t, err)

		cart, err := repo.GetCartByUserIDOrderBySku(ctx, userID)
		require.NoError(t, err)

		require.Empty(t, cart.Items)
	})

	t.Run("get sorted cartItems", func(t *testing.T) {
		t.Parallel()

		repo := NewInMemoryCartRepository(10)
		ctx := context.Background()

		userID := int64(1)
		countItems := 100
		for i := 0; i < countItems; i++ {
			sku, err := rand.Int(rand.Reader, big.NewInt(1000))
			require.NoError(t, err)

			item := &domain.CartItem{
				Sku:   sku.Int64(),
				Name:  "Item 1",
				Count: 1,
				Price: 100,
			}

			_, err = repo.UpsertCartItem(ctx, userID, item)
			require.NoError(t, err)
		}

		cart, err := repo.GetCartByUserIDOrderBySku(ctx, userID)
		require.NoError(t, err)

		require.NotEmpty(t, cart.Items)
		for i := 1; i < len(cart.Items); i++ {
			assert.LessOrEqual(t, cart.Items[i-1].Sku, cart.Items[i].Sku)
		}
	})
}

func BenchmarkUpsertCartItemParallel(b *testing.B) {
	repo := NewInMemoryCartRepository(10)
	ctx := context.Background()
	userID := int64(123)

	cartItems := make([]*domain.CartItem, 100)
	for i := 0; i < 100; i++ {
		sku, err := rand.Int(rand.Reader, big.NewInt(1000000))
		if err != nil {
			panic(err)
		}

		count, err := rand.Int(rand.Reader, big.NewInt(1000))
		if err != nil {
			panic(err)
		}

		cartItems[i] = &domain.CartItem{
			Sku:   sku.Int64(),
			Name:  "any",
			Price: 10,
			Count: uint32(count.Uint64()), // #nosec G115: value is guaranteed to fit in uint32
		}
	}

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			item := cartItems[i%len(cartItems)]
			i++
			_, _ = repo.UpsertCartItem(ctx, userID, item)
		}
	})
}

func BenchmarkCartRepositoryParallel(b *testing.B) {
	repo := NewInMemoryCartRepository(10)
	ctx := context.Background()
	userID := int64(123)

	n := 1000
	cartItems := make([]*domain.CartItem, n)
	for i := 0; i < n; i++ {
		sku, err := rand.Int(rand.Reader, big.NewInt(1000000))
		if err != nil {
			panic(err)
		}

		count, err := rand.Int(rand.Reader, big.NewInt(1000))
		if err != nil {
			panic(err)
		}

		cartItems[i] = &domain.CartItem{
			Sku:   sku.Int64(),
			Name:  "any",
			Price: 10,
			Count: uint32(count.Uint64()), // #nosec G115: value is guaranteed to fit in uint32
		}
	}

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			item := cartItems[i%n]
			skuID := item.Sku
			i++

			_, _ = repo.UpsertCartItem(ctx, userID, item)

			if i%5 == 0 {
				_, _ = repo.GetCartByUserIDOrderBySku(ctx, userID)
			}

			if i%10 == 0 {
				_ = repo.DeleteCartItem(ctx, userID, skuID)
			}

			if i%100 == 0 {
				_ = repo.DeleteCart(ctx, userID)
			}
		}
	})
}
