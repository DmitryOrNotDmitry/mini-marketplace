package handler

import (
	"context"
	"route256/cart/internal/domain"
)

type CartService interface {
	AddCartItem(ctx context.Context, userID int64, newItem *domain.CartItem) (*domain.CartItem, error)
	DeleteCartItem(ctx context.Context, userID, skuID int64) (*domain.CartItem, error)
	ClearCart(ctx context.Context, userID int64) (*domain.Cart, error)
	GetCart(ctx context.Context, userID int64) (*domain.Cart, error)
}

type Server struct {
	cartService CartService
}

func NewServer(cartService CartService) *Server {
	return &Server{
		cartService: cartService,
	}
}
