package service

import (
	"context"
	"fmt"
	"route256/cart/internal/domain"
	"route256/cart/pkg/myerrgroup"
)

// CartRepository описывает методы работы с корзинами в хранилище.
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

// ProductService описывает методы работы с товарами.
type ProductService interface {
	// GetProductBySku возвращает информацию о товаре по SKU.
	GetProductBySku(ctx context.Context, sku int64) (*domain.Product, error)
}

// LomsService описывает методы работы с запасами.
type LomsService interface {
	// GetStockInfo возвращает информацию о запасах товара по SKU.
	GetStockInfo(ctx context.Context, skuID int64) (uint32, error)
}

// CartService содержит бизнес-логику для работы с корзиной.
type CartService struct {
	cartRepository CartRepository
	productService ProductService
	lomsService    LomsService
}

// NewCartService конструктор для CartService.
func NewCartService(repository CartRepository, productService ProductService, lomsService LomsService) *CartService {
	return &CartService{
		cartRepository: repository,
		productService: productService,
		lomsService:    lomsService,
	}
}

// AddCartItem добавляет товар в корзину пользователя, если хватает запасов.
func (s *CartService) AddCartItem(ctx context.Context, userID int64, newItem *domain.CartItem) (*domain.CartItem, error) {
	_, err := s.productService.GetProductBySku(ctx, newItem.Sku)
	if err != nil {
		return nil, fmt.Errorf("productService.GetProductBySku: %w", err)
	}

	productStock, err := s.lomsService.GetStockInfo(ctx, newItem.Sku)
	if err != nil {
		return nil, fmt.Errorf("lomsService.GetStockInfo: %w", err)
	}
	if productStock < newItem.Count {
		return nil, domain.ErrOutOfStock
	}

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

	errGroup, ctx := myerrgroup.WithContext(ctx)

	for _, item := range cart.Items {
		errGroup.Go(func() error {
			product, err := s.productService.GetProductBySku(ctx, item.Sku)
			if err != nil {
				return fmt.Errorf("productService.GetProductBySku: %w", err)
			}

			item.Name = product.Name
			item.Price = product.Price
			return nil
		})
	}

	err = errGroup.Wait()
	if err != nil {
		return nil, fmt.Errorf("errGroup.Wait: %w", err)
	}

	for _, item := range cart.Items {
		cart.TotalPrice += item.Count * item.Price
	}

	return cart, nil
}
