package service

import (
	"context"
	"fmt"
	"route256/cart/internal/domain"
	"route256/loms/pkg/api/orders/v1"
	"route256/loms/pkg/api/stocks/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type LomsServiceGRPC struct {
	stockClient stocks.StockServiceClient
	orderClient orders.OrderServiceClient
}

func NewLomsServiceGRPC(host string, port string) (*LomsServiceGRPC, error) {
	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%s", host, port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("grpc.NewClient: %w", err)
	}

	stockClient := stocks.NewStockServiceClient(conn)
	orderClient := orders.NewOrderServiceClient(conn)

	return &LomsServiceGRPC{
		stockClient: stockClient,
		orderClient: orderClient,
	}, nil
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
