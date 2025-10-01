package postgres

import (
	"context"
	"route256/cart/pkg/metrics"
	"time"

	"github.com/jackc/pgx/v5"
)

type MetricsTracer struct {
}

func NewMetricsTracer() *MetricsTracer {
	return &MetricsTracer{}
}

type ctxKey string

const startTimeKey ctxKey = "queryStartTime"

func (t *MetricsTracer) TraceQueryStart(
	ctx context.Context,
	_ *pgx.Conn,
	_ pgx.TraceQueryStartData,
) context.Context {
	ctx = context.WithValue(ctx, startTimeKey, time.Now())
	return ctx
}

func (t *MetricsTracer) TraceQueryEnd(
	ctx context.Context,
	_ *pgx.Conn,
	data pgx.TraceQueryEndData,
) {
	elapsed := time.Duration(0)
	if v := ctx.Value(startTimeKey); v != nil {
		if start, ok := v.(time.Time); ok {
			elapsed = time.Since(start)
		}
	}

	category := "unexpected"
	switch {
	case data.CommandTag.Select():
		category = metrics.Select
	case data.CommandTag.Update():
		category = metrics.Update
	case data.CommandTag.Delete():
		category = metrics.Delete
	case data.CommandTag.Insert():
		category = metrics.Insert
	}

	metrics.IncDBRequestCount(category)
	metrics.AddDBRequestDurationHist(category, data.Err, elapsed)
}
