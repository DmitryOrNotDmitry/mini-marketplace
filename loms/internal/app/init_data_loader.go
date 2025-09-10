package app

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"route256/cart/pkg/logger"
	"route256/loms/internal/domain"
	"time"
)

type StockCreator interface {
	Create(ctx context.Context, stock *domain.Stock) error
}

// LoadStocks загружает данные из файла в БД
func LoadStocks(filename string, stockCreator StockCreator) error {
	f, err := os.Open(filename) // nolint:gosec
	if err != nil {
		return fmt.Errorf("os.Open: %w", err)
	}
	defer f.Close()

	stocks := []*domain.Stock{}
	err = json.NewDecoder(f).Decode(&stocks)
	if err != nil {
		return fmt.Errorf("json.Decode: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	for _, stock := range stocks {
		errStock := stockCreator.Create(ctx, stock)
		if errStock != nil {
			logger.Error(fmt.Sprintf("Error while try to create stock %+v", stock))
		}
	}

	logger.Info(fmt.Sprintf("Данные из файлы %s загружены успешно", filename))

	return nil
}
