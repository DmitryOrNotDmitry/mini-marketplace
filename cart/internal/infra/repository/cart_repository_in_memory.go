package repository

import (
	"cmp"
	"context"
	"route256/cart/internal/domain"
	"slices"
	"sync"
)

type Storage = map[int64]*domain.Cart

type CartRepositoryInMemory struct {
	storage Storage
	mx      sync.RWMutex
}

func NewInMemoryCartRepository(cap int) *CartRepositoryInMemory {
	return &CartRepositoryInMemory{
		storage: make(Storage, cap),
	}
}

func (r *CartRepositoryInMemory) getOrCreateUserCart(userID int64) (*domain.Cart, error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	cart, ok := r.storage[userID]
	if !ok {
		cart = &domain.Cart{Items: make([]*domain.CartItem, 0, 1)}
		r.storage[userID] = cart
	}

	return cart, nil
}

func (r *CartRepositoryInMemory) UpsertCartItem(_ context.Context, userID int64, newItem *domain.CartItem) (*domain.CartItem, error) {
	cart, err := r.getOrCreateUserCart(userID)
	if err != nil {
		return nil, err
	}

	r.mx.Lock()
	defer r.mx.Unlock()

	for _, item := range cart.Items {
		if item.Sku == newItem.Sku {
			item.Count += newItem.Count
			return item, nil
		}
	}

	cart.Items = append(cart.Items, newItem)

	return newItem, nil
}

func (r *CartRepositoryInMemory) DeleteCartItem(_ context.Context, userID, skuID int64) (*domain.CartItem, error) {
	cart, err := r.getOrCreateUserCart(userID)
	if err != nil {
		return nil, err
	}

	r.mx.Lock()
	defer r.mx.Unlock()

	for i, item := range cart.Items {
		if item.Sku == skuID {
			delItem := item
			cart.Items[i] = cart.Items[len(cart.Items)-1]
			cart.Items = cart.Items[:len(cart.Items)-1]
			return delItem, nil
		}
	}

	return nil, nil
}

func (r *CartRepositoryInMemory) DeleteCart(_ context.Context, userID int64) (*domain.Cart, error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	cart := r.storage[userID]
	delete(r.storage, userID)

	return cart, nil
}

func (r *CartRepositoryInMemory) GetCartByUserIDOrderBySku(_ context.Context, userID int64) (*domain.Cart, error) {
	cart, err := r.getOrCreateUserCart(userID)
	if err != nil {
		return nil, err
	}

	r.mx.RLock()

	cartCopy := *cart
	cartCopy.Items = make([]*domain.CartItem, len(cart.Items))
	copy(cartCopy.Items, cart.Items)

	r.mx.RUnlock()

	slices.SortFunc(cartCopy.Items, func(a, b *domain.CartItem) int {
		return cmp.Compare(a.Sku, b.Sku)
	})

	return &cartCopy, nil
}
