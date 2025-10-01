package interceptor

import (
	"context"
	"route256/cart/pkg/tracer"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// Tracing собирает трейсы с gRPC-запросов.
type Tracing struct {
	tm *tracer.TracerManager
}

// NewTracing создает новый Tracing
func NewTracing(tm *tracer.TracerManager) *Tracing {
	return &Tracing{tm: tm}
}

// Do собирает трейс с запроса.
func (t *Tracing) Do(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	propagator := otel.GetTextMapPropagator()

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.MD{}
	}
	md["Traceparent"] = md["traceparent"]

	ctx = propagator.Extract(ctx, propagation.HeaderCarrier(md))

	ctx, span := t.tm.Tracer.Start(ctx, info.FullMethod)
	defer span.End()

	return handler(ctx, req)
}
