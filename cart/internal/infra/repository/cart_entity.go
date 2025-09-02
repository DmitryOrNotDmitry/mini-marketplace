package repository

import "route256/cart/internal/domain"

type CartEntity struct {
	Items      map[int64]*domain.CartItem
	TotalPrice uint32
}
