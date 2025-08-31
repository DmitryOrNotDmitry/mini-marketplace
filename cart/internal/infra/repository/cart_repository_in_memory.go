package repository

import (
	"cmp"
	"context"
	"route256/cart/internal/domain"
	"slices"
	"sync"
)

type Storage = map[int64]*CartEntity

type CartRepositoryInMemory struct {
	storage Storage
	mx      sync.RWMutex
}

// NewInMemoryCartRepository создает новый репозиторий корзины с in-memory хранилищем.
func NewInMemoryCartRepository(cap int) *CartRepositoryInMemory {
	return &CartRepositoryInMemory{
		storage: make(Storage, cap),
	}
}

func (r *CartRepositoryInMemory) getCartBy(userID int64) (*CartEntity, bool) {
	cart, ok := r.storage[userID]
	return cart, ok
}

func (r *CartRepositoryInMemory) createCartBy(userID int64) *CartEntity {
	cart := &CartEntity{
		Items: make(map[int64]*domain.CartItem),
	}
	r.storage[userID] = cart

	return cart
}

// UpsertCartItem добавляет товар или обновляет количество товара в корзине пользователя в in-memory хранилище.
func (r *CartRepositoryInMemory) UpsertCartItem(_ context.Context, userID int64, newItem *domain.CartItem) (*domain.CartItem, error) {
	r.mx.RLock()
	cart, ok := r.getCartBy(userID)
	r.mx.RUnlock()

	r.mx.Lock()
	defer r.mx.Unlock()

	if !ok {
		cart = r.createCartBy(userID)
	}

	item, ok := cart.Items[newItem.Sku]
	if ok {
		item.Count += newItem.Count
	} else {
		cart.Items[newItem.Sku] = newItem
		item = newItem
	}

	return item, nil
}

// DeleteCartItem удаляет товар из корзины пользователя по SKU из in-memory хранилища.
func (r *CartRepositoryInMemory) DeleteCartItem(_ context.Context, userID, skuID int64) error {
	r.mx.RLock()
	cart, ok := r.getCartBy(userID)
	r.mx.RUnlock()

	r.mx.Lock()
	defer r.mx.Unlock()

	if !ok {
		cart = r.createCartBy(userID)
	}

	delete(cart.Items, skuID)

	return nil
}

// DeleteCart удаляет корзину пользователя из in-memory хранилища.
func (r *CartRepositoryInMemory) DeleteCart(_ context.Context, userID int64) error {
	r.mx.Lock()
	defer r.mx.Unlock()

	delete(r.storage, userID)

	return nil
}

// GetCartByUserIDOrderBySku возвращает корзину пользователя с отсортированными по SKU товарами из in-memory хранилища.
func (r *CartRepositoryInMemory) GetCartByUserIDOrderBySku(_ context.Context, userID int64) (*domain.Cart, error) {
	r.mx.RLock()
	cart, ok := r.getCartBy(userID)
	r.mx.RUnlock()

	if !ok {
		return &domain.Cart{Items: []*domain.CartItem{}}, nil
	}

	r.mx.RLock()

	cartCopy := &domain.Cart{
		Items: make([]*domain.CartItem, 0, len(cart.Items)),
	}
	for _, item := range cart.Items {
		cartCopy.Items = append(cartCopy.Items, item)
	}

	r.mx.RUnlock()

	slices.SortFunc(cartCopy.Items, func(a, b *domain.CartItem) int {
		return cmp.Compare(a.Sku, b.Sku)
	})

	return cartCopy, nil
}
