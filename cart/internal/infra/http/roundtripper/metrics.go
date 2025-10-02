package roundtripper

import (
	"net/http"
	"route256/cart/pkg/metrics"
	"time"
)

// MetricsRoundTripper - http.RoundTripper (middleware) с поддержкой сбора метрик.
type MetricsRoundTripper struct {
	rt http.RoundTripper
}

// NewMetricsRoundTripper создает новый MetricsRoundTripper
func NewMetricsRoundTripper(rt http.RoundTripper) *MetricsRoundTripper {
	return &MetricsRoundTripper{
		rt: rt,
	}
}

func (m *MetricsRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now()

	resp, err := m.rt.RoundTrip(req)

	pattern := req.Pattern
	if pattern == "" {
		pattern = "empty pattern url"
	}

	metrics.IncRequestToExternalCount(pattern)
	metrics.AddRequestToExternalDurationHist(pattern, resp.StatusCode, time.Since(start))

	return resp, err
}
