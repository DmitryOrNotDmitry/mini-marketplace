package repository

import (
	"context"
	"route256/cart/internal/domain"
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

func (r *CartRepositoryInMemory) getOrCreateUserCart(userId int64) (*domain.Cart, error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	cart, ok := r.storage[userId]
	if !ok {
		cart = &domain.Cart{Items: make([]*domain.CartItem, 0, 1)}
		r.storage[userId] = cart
	}

	return cart, nil
}

func (r *CartRepositoryInMemory) AddCartItem(ctx context.Context, userId int64, newItem *domain.CartItem) (*domain.CartItem, error) {
	cart, err := r.getOrCreateUserCart(userId)
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

func (r *CartRepositoryInMemory) DeleteCartItem(ctx context.Context, userId, skuId int64) (*domain.CartItem, error) {
	cart, err := r.getOrCreateUserCart(userId)
	if err != nil {
		return nil, err
	}

	r.mx.Lock()
	defer r.mx.Unlock()

	for i, item := range cart.Items {
		if item.Sku == skuId {
			delItem := item
			cart.Items[i] = cart.Items[len(cart.Items)-1]
			cart.Items = cart.Items[:len(cart.Items)-1]
			return delItem, nil
		}
	}

	return nil, nil
}

func (r *CartRepositoryInMemory) ClearCart(ctx context.Context, userId int64) (*domain.Cart, error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	cart := r.storage[userId]
	delete(r.storage, userId)

	return cart, nil
}

func (r *CartRepositoryInMemory) GetCart(ctx context.Context, userId int64) (*domain.Cart, error) {
	cart, err := r.getOrCreateUserCart(userId)
	if err != nil {
		return nil, err
	}

	cartCopy := *cart
	cartCopy.Items = make([]*domain.CartItem, len(cart.Items))
	copy(cartCopy.Items, cart.Items)

	return &cartCopy, nil
}
