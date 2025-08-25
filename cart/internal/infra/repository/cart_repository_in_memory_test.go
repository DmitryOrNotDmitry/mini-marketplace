package repository

import (
	"context"
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
				r.UpsertCartItem(context.Background(), 1, &domain.CartItem{Sku: 100, Count: 2, Name: "Test", Price: 100})
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
				r.UpsertCartItem(context.Background(), 1, &domain.CartItem{Sku: 200, Count: 1})
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
			got, err := tt.repo.DeleteCartItem(tt.args.ctx, tt.args.userID, tt.args.skuID)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteCartItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteCartItem() = %v, want %v", got, tt.want)
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
				r.UpsertCartItem(context.Background(), 1, &domain.CartItem{Sku: 300, Count: 1})
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
			got, err := tt.repo.DeleteCart(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteCart() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteCart() = %v, want %v", got, tt.want)
			}

			if tt.want != nil {
				cartAfter, _ := tt.repo.GetCartByUserIdOrderBySku(tt.args.ctx, tt.args.userID)
				if len(cartAfter.Items) != 0 {
					t.Errorf("cart was not deleted, got %v", cartAfter.Items)
				}
			}
		})
	}
}

func TestCartRepositoryInMemory_GetCartByUserIdOrderBySku(t *testing.T) {
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
				r.UpsertCartItem(context.Background(), 1, &domain.CartItem{Sku: 300, Count: 1})
				r.UpsertCartItem(context.Background(), 1, &domain.CartItem{Sku: 100, Count: 2})
				r.UpsertCartItem(context.Background(), 1, &domain.CartItem{Sku: 200, Count: 3})
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
			got, err := tt.repo.GetCartByUserIdOrderBySku(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCartByUserIdOrderBySku() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCartByUserIdOrderBySku() = %v, want %v", got, tt.want)
			}
		})
	}
}
