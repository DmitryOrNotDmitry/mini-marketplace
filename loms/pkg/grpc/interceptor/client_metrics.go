package interceptor

import (
	"context"
	"route256/cart/pkg/metrics"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// ClientMetrics собирает метрики с gRPC-запросы с клиентской стороны.
func ClientMetrics(
	ctx context.Context,
	method string,
	req, reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	start := time.Now()
	err := invoker(ctx, method, req, reply, cc, opts...)

	var codeInt = -1
	st, ok := status.FromError(err)
	if ok {
		codeInt = int(st.Code())
	}

	metrics.AddRequestToExternalDurationHist(method, codeInt, time.Since(start))
	metrics.IncRequestToExternalCount(method)

	return err
}
