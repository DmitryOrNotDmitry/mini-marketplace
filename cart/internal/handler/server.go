package handler

import (
	"context"
	"route256/cart/internal/domain"
)

type CartService interface {
	AddCartItem(ctx context.Context, userId int64, newItem *domain.CartItem) (*domain.CartItem, error)
	DeleteCartItem(ctx context.Context, userId, skuId int64) (*domain.CartItem, error)
}

type Server struct {
	cartService CartService
}

func NewServer(cartService CartService) *Server {
	return &Server{
		cartService: cartService,
	}
}
