package app

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"time"

	"route256/cart/pkg/logger"
	"route256/loms/internal/domain"
)

type StockCreator interface {
	Create(ctx context.Context, stock *domain.Stock) error
}

//go:embed init_data/stock-data.json
var stocksData []byte

// LoadStocks загружает данные из файла в БД
func LoadStocks(stockCreator StockCreator, timeout time.Duration) error {
	stocks := []*domain.Stock{}
	err := json.NewDecoder(bytes.NewReader(stocksData)).Decode(&stocks)
	if err != nil {
		return fmt.Errorf("json.Decode: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	for _, stock := range stocks {
		errStock := stockCreator.Create(ctx, stock)
		if errStock != nil {
			logger.Error(fmt.Sprintf("Error while try to create stock %+v", stock))
		}
	}

	logger.Info(fmt.Sprintf("Данные из файлы %s загружены успешно", "init_data/stock-data.json"))

	return nil
}
