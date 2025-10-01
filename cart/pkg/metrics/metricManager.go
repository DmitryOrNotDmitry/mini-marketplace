package metrics

import (
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type DBRequestCategory = string

const (
	Select DBRequestCategory = "select"
	Update DBRequestCategory = "update"
	Insert DBRequestCategory = "insert"
	Delete DBRequestCategory = "delete"
)

var (
	requestCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "handler_request_total_counter",
		Help: "Total count of request",
	}, []string{"handler"})

	requestDurationHistogram = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "handler_request_duration_historgram",
		Help:    "Total duration of handler processing",
		Buckets: prometheus.DefBuckets,
	}, []string{"handler", "status"})

	requestToExternalCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "request_to_external_services_total_counter",
		Help: "Total count of request to external services",
	}, []string{"resource"})

	requestToExternalDurationHistogram = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "request_to_external_services_duration_historgram",
		Help:    "Total duration of request to external services",
		Buckets: prometheus.DefBuckets,
	}, []string{"resource", "status"})

	requestToDBCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "database_request_total_counter",
		Help: "Total count of request",
	}, []string{"category"})

	requestToDBDurationHistogram = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "database_request_duration_historgram",
		Help:    "Total duration of request to database",
		Buckets: prometheus.DefBuckets,
	}, []string{"category", "error"})

	repositorySizeGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "repo_size_gauge",
		Help: "Size of repo",
	}, []string{"object"})
)

func IncRequestCount(handler string) {
	requestCounter.WithLabelValues(handler).Inc()
}

func AddRequestDurationHist(handler string, status int, duration time.Duration) {
	requestDurationHistogram.
		WithLabelValues(handler, fmt.Sprint(status)).
		Observe(float64(duration.Seconds()))
}

func IncRequestToExternalCount(resource string) {
	requestToExternalCounter.WithLabelValues(resource).Inc()
}

func AddRequestToExternalDurationHist(resource string, status int, duration time.Duration) {
	requestToExternalDurationHistogram.
		WithLabelValues(resource, fmt.Sprint(status)).
		Observe(float64(duration.Seconds()))
}

func IncDBRequestCount(category DBRequestCategory) {
	requestToDBCounter.WithLabelValues(category).Inc()
}

func AddDBRequestDurationHist(category DBRequestCategory, err error, duration time.Duration) {
	errData := "no"
	if err != nil {
		errData = "yes"
	}

	requestToDBDurationHistogram.
		WithLabelValues(category, errData).
		Observe(float64(duration.Seconds()))
}

func StoreRepositorySize(objectsName string, size float64) {
	repositorySizeGauge.WithLabelValues(objectsName).Set(size)
}
