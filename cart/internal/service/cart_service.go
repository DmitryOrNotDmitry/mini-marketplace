package service

import (
	"context"
	"fmt"
	"route256/cart/internal/domain"
)

type CartRepository interface {
	GetCartByUserIDOrderBySku(ctx context.Context, userID int64) (*domain.Cart, error)
	DeleteCart(ctx context.Context, userID int64) (*domain.Cart, error)

	UpsertCartItem(ctx context.Context, userID int64, newItem *domain.CartItem) (*domain.CartItem, error)
	DeleteCartItem(ctx context.Context, userID, skuID int64) (*domain.CartItem, error)
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

func (s *CartService) AddCartItem(ctx context.Context, userID int64, newItem *domain.CartItem) (*domain.CartItem, error) {
	product, err := s.productService.GetProductBySku(ctx, newItem.Sku)
	if err != nil {
		return nil, fmt.Errorf("productService.GetProductBySku: %w", err)
	}

	newItem.Price = uint32(product.Price)
	newItem.Name = product.Name

	addedCartItem, err := s.cartRepository.UpsertCartItem(ctx, userID, newItem)
	if err != nil {
		return nil, fmt.Errorf("cartRepository.UpsertCartItem: %w", err)
	}

	return addedCartItem, nil
}

func (s *CartService) DeleteCartItem(ctx context.Context, userID, skuID int64) (*domain.CartItem, error) {
	deletedCartItem, err := s.cartRepository.DeleteCartItem(ctx, userID, skuID)
	if err != nil {
		return nil, fmt.Errorf("cartRepository.DeleteCartItem: %w", err)
	}

	return deletedCartItem, nil
}

func (s *CartService) ClearCart(ctx context.Context, userID int64) (*domain.Cart, error) {
	deletedCart, err := s.cartRepository.DeleteCart(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("cartRepository.DeleteCart: %w", err)
	}

	return deletedCart, nil
}

func (s *CartService) GetCart(ctx context.Context, userID int64) (*domain.Cart, error) {
	cart, err := s.cartRepository.GetCartByUserIDOrderBySku(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("cartRepository.GetCartByUserIDOrderBySku: %w", err)
	}

	for _, item := range cart.Items {
		cart.TotalPrice += item.Count * item.Price
	}

	return cart, nil
}
