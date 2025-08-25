package service

import (
	"context"
	"errors"
	"route256/cart/internal/domain"
	"testing"

	"github.com/stretchr/testify/mock"
)

type MockCartRepository struct {
	mock.Mock
}

func (m *MockCartRepository) UpsertCartItem(ctx context.Context, userId int64, newItem *domain.CartItem) (*domain.CartItem, error) {
	args := m.Called(ctx, userId, newItem)
	return args.Get(0).(*domain.CartItem), args.Error(1)
}

func (m *MockCartRepository) DeleteCartItem(ctx context.Context, userId, skuId int64) (*domain.CartItem, error) {
	args := m.Called(ctx, userId, skuId)
	return args.Get(0).(*domain.CartItem), args.Error(1)
}

func (m *MockCartRepository) DeleteCart(ctx context.Context, userId int64) (*domain.Cart, error) {
	args := m.Called(ctx, userId)
	return args.Get(0).(*domain.Cart), args.Error(1)
}

func (m *MockCartRepository) GetCartByUserIdOrderBySku(ctx context.Context, userId int64) (*domain.Cart, error) {
	args := m.Called(ctx, userId)
	return args.Get(0).(*domain.Cart), args.Error(1)
}

type MockProductService struct {
	mock.Mock
}

func (m *MockProductService) GetProductBySku(ctx context.Context, sku int64) (*domain.Product, error) {
	args := m.Called(ctx, sku)
	return args.Get(0).(*domain.Product), args.Error(1)
}

func TestCartService_AddCartItem(t *testing.T) {
	mockRepo := new(MockCartRepository)
	mockProduct := new(MockProductService)
	service := NewCartService(mockRepo, mockProduct)

	newItem := &domain.CartItem{Sku: 1, Count: 2}
	expectedProduct := &domain.Product{Name: "TestProd", Price: 100, Sku: 1}
	expectedCartItem := &domain.CartItem{Sku: 1, Count: 2, Name: "TestProd", Price: 100}

	mockProduct.On("GetProductBySku", mock.Anything, int64(1)).Return(expectedProduct, nil)
	mockRepo.On("UpsertCartItem", mock.Anything, int64(123), mock.Anything).Return(expectedCartItem, nil)

	got, err := service.AddCartItem(context.Background(), 123, newItem)
	if err != nil {
		t.Fatal(err)
	}

	if got.Name != expectedCartItem.Name || got.Price != expectedCartItem.Price {
		t.Errorf("unexpected result: %+v", got)
	}

	mockRepo.AssertExpectations(t)
	mockProduct.AssertExpectations(t)
}

func TestCartService_GetCart(t *testing.T) {
	mockRepo := new(MockCartRepository)
	service := NewCartService(mockRepo, nil)

	cart := &domain.Cart{
		Items: []*domain.CartItem{
			{Sku: 1, Count: 2, Price: 100},
			{Sku: 2, Count: 1, Price: 50},
		},
	}
	mockRepo.On("GetCartByUserIdOrderBySku", mock.Anything, int64(123)).Return(cart, nil)

	got, err := service.GetCart(context.Background(), 123)
	if err != nil {
		t.Fatal(err)
	}

	expectedTotal := uint32(2*100 + 1*50)
	if got.TotalPrice != expectedTotal {
		t.Errorf("expected total %d, got %d", expectedTotal, got.TotalPrice)
	}

	mockRepo.AssertExpectations(t)
}

func TestCartService_AddCartItem_ProductError(t *testing.T) {
	mockProduct := new(MockProductService)
	service := NewCartService(nil, mockProduct)

	var nilProduct *domain.Product
	mockProduct.On("GetProductBySku", mock.Anything, int64(1)).Return(nilProduct, errors.New("product not found"))

	_, err := service.AddCartItem(context.Background(), 123, &domain.CartItem{Sku: 1, Count: 1})
	if err == nil {
		t.Errorf("expected error but got nil")
	}

	mockProduct.AssertExpectations(t)
}
