package middleware

import (
	"net/http"
	"route256/cart/pkg/tracer"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

// Tracing - middleware для сбора трайсинга с HTTP-запросов.
type Tracing struct {
	h  http.Handler
	tm *tracer.TracerManager
}

// NewTracing создает middleware для сбора метрик HTTP-запросов.
func NewTracing(h http.Handler, tm *tracer.TracerManager) http.Handler {
	return &Tracing{h: h, tm: tm}
}

func (t *Tracing) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	propagator := otel.GetTextMapPropagator()
	ctx := propagator.Extract(r.Context(), propagation.HeaderCarrier(r.Header))

	ctx, span := t.tm.Tracer.Start(ctx, "none")
	defer func() {
		span.SetName(r.Pattern)
		span.End()
	}()

	r = r.WithContext(ctx)

	t.h.ServeHTTP(w, r)
}
