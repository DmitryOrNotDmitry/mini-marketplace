package handler

import (
	"context"
	"errors"
	"route256/loms/internal/domain"
	"route256/loms/pkg/api/stocks/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type StockService interface {
	GetAvailableCount(ctx context.Context, skuID int64) (uint32, error)
}

type StockServerGRPC struct {
	stockService StockService
	stocks.UnimplementedStockServiceServer
}

// NewStockServerGRPC создает новый экземпляр gRPC-сервера StockServerGRPC.
func NewStockServerGRPC(stockService StockService) *StockServerGRPC {
	return &StockServerGRPC{
		stockService: stockService,
	}
}

// StockInfo обрабатывает gRPC-запрос на получение информации о товаре по его SKU.
func (ss *StockServerGRPC) StockInfo(ctx context.Context, req *stocks.StockInfoRequest) (*stocks.StockInfoResponse, error) {
	count, err := ss.stockService.GetAvailableCount(ctx, req.SkuId)
	if err != nil {
		if errors.Is(err, domain.ErrItemStockNotExist) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &stocks.StockInfoResponse{Count: count}, nil
}
