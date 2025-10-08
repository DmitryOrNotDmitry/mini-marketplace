package interceptor

import (
	"context"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// ClientTracing реализует gRPC клиентский интерцептор для передачи контекста трассировки.
func ClientTracing(
	ctx context.Context,
	method string,
	req, reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	propagator := otel.GetTextMapPropagator()

	md := metadata.New(nil)
	propagator.Inject(ctx, propagation.HeaderCarrier(md))

	norm := metadata.MD{}
	for k, v := range md {
		norm[strings.ToLower(k)] = v
	}

	ctx = metadata.NewOutgoingContext(ctx, norm)

	return invoker(ctx, method, req, reply, cc, opts...)
}
