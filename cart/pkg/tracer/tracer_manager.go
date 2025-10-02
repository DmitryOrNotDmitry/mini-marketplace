package tracer

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	t "go.opentelemetry.io/otel/trace"
)

// Manager управляет трассировщиком и провайдером трассировки.
type Manager struct {
	Tracer   t.Tracer
	provider *trace.TracerProvider
}

// NewTracerManager создает новый экземпляр Manager для трассировки.
func NewTracerManager(ctx context.Context, jaegerUrl, serviceName, enviroment string) (*Manager, error) {
	exp, err := otlptracehttp.New(ctx, otlptracehttp.WithEndpointURL(jaegerUrl))
	if err != nil {
		return nil, err
	}
	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(serviceName),
			semconv.DeploymentEnvironmentName(enviroment),
		),
	)
	if err != nil {
		return nil, err
	}
	tracerProvider := trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(r),
	)

	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
		),
	)

	tracer := otel.GetTracerProvider().Tracer(serviceName)
	return &Manager{
		Tracer:   tracer,
		provider: tracerProvider,
	}, nil
}

// Stop завершает работу провайдера трассировки.
func (t *Manager) Stop(ctx context.Context) error {
	return t.provider.Shutdown(ctx)
}
