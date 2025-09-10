package service

import (
	"context"
	"route256/cart/internal/domain"
	"route256/loms/pkg/api/orders/v1"
	"route256/loms/pkg/api/stocks/v1"
)

type LomsServiceGRPC struct {
	stockClient stocks.StockServiceClient
	orderClient orders.OrderServiceClient
}

func NewLomsServiceGRPC(stockClient stocks.StockServiceClient, orderClient orders.OrderServiceClient) *LomsServiceGRPC {
	return &LomsServiceGRPC{
		stockClient: stockClient,
		orderClient: orderClient,
	}
}

func (ls *LomsServiceGRPC) GetStockInfo(ctx context.Context, skuID int64) (uint32, error) {
	resp, err := ls.stockClient.StockInfo(ctx, &stocks.StockInfoRequest{
		SkuId: skuID,
	})
	if err != nil {
		return 0, err
	}

	return resp.Count, nil
}

func (ls *LomsServiceGRPC) OrderCreate(ctx context.Context, userID int64, cart *domain.Cart) (int64, error) {
	req := &orders.OrderCreateRequest{
		UserId: userID,
		Items:  make([]*orders.ItemInfo, 0, len(cart.Items)),
	}
	for _, item := range cart.Items {
		req.Items = append(req.Items, &orders.ItemInfo{
			SkuId: item.Sku,
			Count: item.Count,
		})
	}

	resp, err := ls.orderClient.OrderCreate(ctx, req)
	if err != nil {
		return 0, err
	}

	return resp.OrderId, nil
}
