package service

import (
	"context"
	"fmt"
	"route256/cart/internal/domain"
)

type CartRepository interface {
	// GetCartByUserIDOrderBySku возвращает корзину пользователя с отсортированными по SKU товарами.
	GetCartByUserIDOrderBySku(ctx context.Context, userID int64) (*domain.Cart, error)
	// DeleteCart удаляет корзину пользователя.
	DeleteCart(ctx context.Context, userID int64) error

	//  UpsertCartItem добавляет товар или обновляет количество товара в корзине пользователя.
	UpsertCartItem(ctx context.Context, userID int64, newItem *domain.CartItem) (*domain.CartItem, error)
	//  DeleteCartItem удаляет товар из корзины пользователя по SKU.
	DeleteCartItem(ctx context.Context, userID, skuID int64) error
}

type ProductService interface {
	// GetProductBySku возвращает информацию о товаре по SKU.
	GetProductBySku(ctx context.Context, sku int64) (*domain.Product, error)
}

type CartService struct {
	cartRepository CartRepository
	productService ProductService
}

// NewCartService конструктор для CartService.
func NewCartService(repository CartRepository, productService ProductService) *CartService {
	return &CartService{cartRepository: repository, productService: productService}
}

// AddCartItem добавляет товар в корзину пользователя.
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

// DeleteCartItem удаляет товар из корзины пользователя.
func (s *CartService) DeleteCartItem(ctx context.Context, userID, skuID int64) error {
	err := s.cartRepository.DeleteCartItem(ctx, userID, skuID)
	if err != nil {
		return fmt.Errorf("cartRepository.DeleteCartItem: %w", err)
	}

	return nil
}

// ClearCart очищает корзину пользователя.
func (s *CartService) ClearCart(ctx context.Context, userID int64) error {
	err := s.cartRepository.DeleteCart(ctx, userID)
	if err != nil {
		return fmt.Errorf("cartRepository.DeleteCart: %w", err)
	}

	return nil
}

// GetCart возвращает содержимое корзины пользователя.
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
