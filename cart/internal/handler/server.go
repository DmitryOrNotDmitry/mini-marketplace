package handler

import (
	"context"
	"route256/cart/internal/domain"
)

// CartService для работы с корзиной.
type CartService interface {
	// Добавляет товар в корзину пользователя
	AddCartItem(ctx context.Context, userID int64, newItem *domain.CartItem) (*domain.CartItem, error)
	// Удаляет товар из корзины пользователя
	DeleteCartItem(ctx context.Context, userID, skuID int64) error
	// Очищает корзину пользователя
	ClearCart(ctx context.Context, userID int64) error
	// Возвращает содержимое корзины пользователя
	GetCart(ctx context.Context, userID int64) (*domain.Cart, error)
}

// Server реализует HTTP-обработчики для работы с корзиной.
type Server struct {
	cartService CartService
}

// NewServer конструктор для Server.
func NewServer(cartService CartService) *Server {
	return &Server{
		cartService: cartService,
	}
}
