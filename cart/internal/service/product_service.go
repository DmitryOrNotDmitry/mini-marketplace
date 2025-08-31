package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"route256/cart/internal/domain"
	"time"
)

var ErrNotOk = errors.New("status not ok")

type ProductServiceHTTP struct {
	httpClient http.Client
	token      string
	address    string
}

// NewProductServiceHTTP конструктор для ProductServiceHTTP.
func NewProductServiceHTTP(httpClient http.Client, token string, address string) *ProductServiceHTTP {
	return &ProductServiceHTTP{
		httpClient: httpClient,
		token:      token,
		address:    address,
	}
}

type GetProductResponse struct {
	Name  string `json:"name"`
	Price int32  `json:"price"`
	Sku   int64  `json:"sku"`
}

// GetProductBySku возвращает информацию о товаре по SKU из внешнего сервиса.
func (s *ProductServiceHTTP) GetProductBySku(ctx context.Context, sku int64) (*domain.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/product/%d", s.address, sku),
		http.NoBody,
	)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequestWithContext: %w", err)
	}

	req.Header.Add("X-API-KEY", s.token)

	response, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("httpClient.Do: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusNotFound {
		return nil, domain.ErrProductNotFound
	}

	if response.StatusCode != http.StatusOK {
		return nil, ErrNotOk
	}

	resp := &GetProductResponse{}
	if err := json.NewDecoder(response.Body).Decode(resp); err != nil {
		return nil, fmt.Errorf("json.NewDecoder: %w", err)
	}

	return &domain.Product{
		Name:  resp.Name,
		Price: uint32(resp.Price),
		Sku:   resp.Sku,
	}, nil
}
