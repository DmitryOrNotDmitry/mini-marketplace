package service

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"reflect"
	"route256/cart/internal/domain"
	"testing"
)

type mockRoundTripper struct {
	resp *http.Response
	err  error
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.resp, m.err
}

func TestProductServiceHttp_GetProductBySku(t *testing.T) {
	// Мокнутый ответ для 200
	successResp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBufferString(`{"name":"TestProduct","price":100,"sku":12345}`)),
	}
	// Мокнутый ответ для 404
	notFoundResp := &http.Response{
		StatusCode: http.StatusNotFound,
		Body:       io.NopCloser(bytes.NewBufferString(``)),
	}

	tests := []struct {
		name string
		s    *ProductServiceHttp
		args struct {
			ctx context.Context
			sku int64
		}
		want    *domain.Product
		wantErr bool
	}{
		{
			name: "successful response",
			s: &ProductServiceHttp{
				httpClient: http.Client{Transport: &mockRoundTripper{resp: successResp}},
				token:      "token",
				address:    "http://localhost",
			},
			args: struct {
				ctx context.Context
				sku int64
			}{ctx: context.Background(), sku: 12345},
			want:    &domain.Product{Name: "TestProduct", Price: 100, Sku: 12345},
			wantErr: false,
		},
		{
			name: "product not found",
			s: &ProductServiceHttp{
				httpClient: http.Client{Transport: &mockRoundTripper{resp: notFoundResp}},
				token:      "token",
				address:    "http://localhost",
			},
			args: struct {
				ctx context.Context
				sku int64
			}{ctx: context.Background(), sku: 99999},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.GetProductBySku(tt.args.ctx, tt.args.sku)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProductServiceHttp.GetProductBySku() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProductServiceHttp.GetProductBySku() = %v, want %v", got, tt.want)
			}
		})
	}
}
