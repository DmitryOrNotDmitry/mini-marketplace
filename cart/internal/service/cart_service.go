package service

import (
	"context"
	"fmt"
	"route256/cart/internal/domain"
)

type CartRepository interface {
	AddCartItem(ctx context.Context, userId int64, newItem *domain.CartItem) (*domain.CartItem, error)
}

type ProductService interface {
	GetProductBySku(ctx context.Context, sku int64) (*domain.Product, error)
}

type CartService struct {
	cartRepository CartRepository
	productService ProductService
}

func NewCartService(repository CartRepository, productService ProductService) *CartService {
	return &CartService{cartRepository: repository, productService: productService}
}

func (s *CartService) AddCartItem(ctx context.Context, userId int64, newItem *domain.CartItem) (*domain.CartItem, error) {
	product, err := s.productService.GetProductBySku(ctx, newItem.Sku)
	if err != nil {
		return nil, fmt.Errorf("productService.GetProductBySku: %w", err)
	}

	newItem.Price = uint32(product.Price)
	newItem.Name = product.Name

	addedCartItem, err := s.cartRepository.AddCartItem(ctx, userId, newItem)
	if err != nil {
		return nil, fmt.Errorf("cartRepository.AddCartItem: %w", err)
	}

	return addedCartItem, nil
}
