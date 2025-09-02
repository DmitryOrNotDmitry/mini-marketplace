package repository

import (
	"context"
	"crypto/rand"
	"math/big"
	"reflect"
	"route256/cart/internal/domain"
	"testing"
)

func TestCartRepositoryInMemory_UpsertCartItem(t *testing.T) {
	type args struct {
		ctx     context.Context
		userID  int64
		newItem *domain.CartItem
	}
	tests := []struct {
		name    string
		repo    *CartRepositoryInMemory
		args    args
		want    *domain.CartItem
		wantErr bool
	}{
		{
			name: "add new item",
			repo: NewInMemoryCartRepository(10),
			args: args{ctx: context.Background(), userID: 1, newItem: &domain.CartItem{Sku: 100, Count: 2, Name: "Test", Price: 100}},
			want: &domain.CartItem{Sku: 100, Count: 2, Name: "Test", Price: 100},
		},
		{
			name: "increase count for existing SKU",
			repo: func() *CartRepositoryInMemory {
				r := NewInMemoryCartRepository(10)
				_, err := r.UpsertCartItem(context.Background(), 1, &domain.CartItem{Sku: 100, Count: 2, Name: "Test", Price: 100})
				if err != nil {
					t.Fatalf("r.UpsertCartItem returns error - %s", err.Error())
					return nil
				}
				return r
			}(),
			args: args{ctx: context.Background(), userID: 1, newItem: &domain.CartItem{Sku: 100, Count: 3, Name: "Test", Price: 100}},
			want: &domain.CartItem{Sku: 100, Count: 5, Name: "Test", Price: 100},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.repo.UpsertCartItem(tt.args.ctx, tt.args.userID, tt.args.newItem)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpsertCartItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpsertCartItem() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCartRepositoryInMemory_DeleteCartItem(t *testing.T) {
	type args struct {
		ctx    context.Context
		userID int64
		skuID  int64
	}
	tests := []struct {
		name    string
		repo    *CartRepositoryInMemory
		args    args
		want    *domain.CartItem
		wantErr bool
	}{
		{
			name: "delete existing item",
			repo: func() *CartRepositoryInMemory {
				r := NewInMemoryCartRepository(10)
				_, err := r.UpsertCartItem(context.Background(), 1, &domain.CartItem{Sku: 200, Count: 1})
				if err != nil {
					t.Fatalf("r.UpsertCartItem returns error - %s", err.Error())
					return nil
				}
				return r
			}(),
			args: args{ctx: context.Background(), userID: 1, skuID: 200},
			want: &domain.CartItem{Sku: 200, Count: 1},
		},
		{
			name: "delete non-existing item",
			repo: NewInMemoryCartRepository(10),
			args: args{ctx: context.Background(), userID: 1, skuID: 999},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.repo.DeleteCartItem(tt.args.ctx, tt.args.userID, tt.args.skuID)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteCartItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestCartRepositoryInMemory_DeleteCart(t *testing.T) {
	type args struct {
		ctx    context.Context
		userID int64
	}
	tests := []struct {
		name    string
		repo    *CartRepositoryInMemory
		args    args
		want    *domain.Cart
		wantErr bool
	}{
		{
			name: "delete existing cart",
			repo: func() *CartRepositoryInMemory {
				r := NewInMemoryCartRepository(10)
				_, err := r.UpsertCartItem(context.Background(), 1, &domain.CartItem{Sku: 300, Count: 1})
				if err != nil {
					t.Fatalf("r.UpsertCartItem returns error - %s", err.Error())
					return nil
				}
				return r
			}(),
			args: args{ctx: context.Background(), userID: 1},
			want: &domain.Cart{
				Items: []*domain.CartItem{
					{Sku: 300, Count: 1},
				},
			},
		},
		{
			name: "delete non-existing cart",
			repo: NewInMemoryCartRepository(10),
			args: args{ctx: context.Background(), userID: 2},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.repo.DeleteCart(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteCart() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.want != nil {
				cartAfter, _ := tt.repo.GetCartByUserIDOrderBySku(tt.args.ctx, tt.args.userID)
				if len(cartAfter.Items) != 0 {
					t.Errorf("cart was not deleted, got %v", cartAfter.Items)
				}
			}
		})
	}
}

func TestCartRepositoryInMemory_GetCartByUserIDOrderBySku(t *testing.T) {
	type args struct {
		ctx    context.Context
		userID int64
	}
	tests := []struct {
		name    string
		repo    *CartRepositoryInMemory
		args    args
		want    *domain.Cart
		wantErr bool
	}{
		{
			name: "get cart sorted by SKU",
			repo: func() *CartRepositoryInMemory {
				r := NewInMemoryCartRepository(10)
				_, err1 := r.UpsertCartItem(context.Background(), 1, &domain.CartItem{Sku: 300, Count: 1})
				_, err2 := r.UpsertCartItem(context.Background(), 1, &domain.CartItem{Sku: 100, Count: 2})
				_, err3 := r.UpsertCartItem(context.Background(), 1, &domain.CartItem{Sku: 200, Count: 3})
				for _, err := range []error{err1, err2, err3} {
					if err != nil {
						t.Fatalf("r.UpsertCartItem returns error - %s", err.Error())
						return nil
					}
				}
				return r
			}(),
			args: args{ctx: context.Background(), userID: 1},
			want: &domain.Cart{
				Items: []*domain.CartItem{
					{Sku: 100, Count: 2},
					{Sku: 200, Count: 3},
					{Sku: 300, Count: 1},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.repo.GetCartByUserIDOrderBySku(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCartByUserIDOrderBySku() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCartByUserIDOrderBySku() = %v, want %v", got, tt.want)
			}
		})
	}
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
