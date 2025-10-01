package middleware

import (
	"net/http"
	"route256/cart/pkg/metrics"
	"time"
)

// MetricsMiddleware - middleware для сбора метрик HTTP-запросов.
type MetricsMiddleware struct {
	h http.Handler
}

// NewMetricsMiddleware создает middleware для сбора метрик HTTP-запросов.
func NewMetricsMiddleware(h http.Handler) http.Handler {
	return &MetricsMiddleware{h: h}
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.status = code
	rec.ResponseWriter.WriteHeader(code)
}

func (m *MetricsMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	wDecorator := &statusRecorder{
		ResponseWriter: w,
	}

	start := time.Now()
	defer func() {
		pattern := r.Pattern
		if pattern == "" {
			pattern = r.URL.Path
		}

		metrics.AddRequestDurationHist(pattern, wDecorator.status, time.Since(start))
		metrics.IncRequestCount(pattern)
	}()

	m.h.ServeHTTP(wDecorator, r)
}
