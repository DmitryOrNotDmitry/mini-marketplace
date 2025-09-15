package repository

import "route256/cart/internal/domain"

// CartEntity хранит данные о корзине.
// Используется только в репозитории для быстрого доступа к товарам в корзине.
type CartEntity struct {
	Items      map[int64]*domain.CartItem
	TotalPrice uint32
}
