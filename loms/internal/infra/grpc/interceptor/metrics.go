package interceptor

import (
	"context"
	"route256/cart/pkg/metrics"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// Metrics собирает метрики с gRPC-запросов.
func Metrics(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	start := time.Now()

	resp, err := handler(ctx, req)

	st := status.Convert(err)
	metrics.AddRequestDurationHist(info.FullMethod, int(st.Code()), time.Since(start))
	metrics.IncRequestCount(info.FullMethod)

	return resp, err
}
