package middleware

import (
	"fmt"
	"net/http"
	"time"

	"route256/cart/pkg/logger"
)

type LoggerMiddleware struct {
	h http.Handler
}

// NewLoggerMiddleware создает middleware для логирования HTTP-запросов.
func NewLoggerMiddleware(h http.Handler) http.Handler {
	return &LoggerMiddleware{h: h}
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (m *LoggerMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sw := &statusWriter{ResponseWriter: w, status: -1}

	defer func(now time.Time) {
		logger.Info(fmt.Sprintf("%s %s -> %d (%s)",
			r.Method, r.URL.String(), sw.status, time.Since(now)))
	}(time.Now())

	m.h.ServeHTTP(sw, r)
}
