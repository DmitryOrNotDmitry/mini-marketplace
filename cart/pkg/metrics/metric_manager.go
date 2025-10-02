package metrics

import (
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// DBRequestCategory определяет тип запроса к базе данных.
type DBRequestCategory = string

const (
	// Select категория для SELECT-запросов.
	Select DBRequestCategory = "select"
	// Update категория для UPDATE-запросов.
	Update DBRequestCategory = "update"
	// Insert категория для INSERT-запросов.
	Insert DBRequestCategory = "insert"
	// Delete категория для DELETE-запросов.
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

// IncRequestCount увеличивает метрику счетчика запросов для указанного обработчика.
func IncRequestCount(handler string) {
	requestCounter.WithLabelValues(handler).Inc()
}

// AddRequestDurationHist добавляет значение в гистограмму длительности обработки запроса.
func AddRequestDurationHist(handler string, status int, duration time.Duration) {
	requestDurationHistogram.
		WithLabelValues(handler, fmt.Sprint(status)).
		Observe(float64(duration.Seconds()))
}

// IncRequestToExternalCount увеличивает метрику счетчика запросов к внешнему ресурсу.
func IncRequestToExternalCount(resource string) {
	requestToExternalCounter.WithLabelValues(resource).Inc()
}

// AddRequestToExternalDurationHist добавляет значение в гистограмму длительности запроса к внешнему ресурсу.
func AddRequestToExternalDurationHist(resource string, status int, duration time.Duration) {
	requestToExternalDurationHistogram.
		WithLabelValues(resource, fmt.Sprint(status)).
		Observe(float64(duration.Seconds()))
}

// IncDBRequestCount увеличивает метрику счетчика запросов к базе данных по категории.
func IncDBRequestCount(category DBRequestCategory) {
	requestToDBCounter.WithLabelValues(category).Inc()
}

// AddDBRequestDurationHist добавляет значение в гистограмму длительности запроса к базе данных.
func AddDBRequestDurationHist(category DBRequestCategory, err error, duration time.Duration) {
	errData := "no"
	if err != nil {
		errData = "yes"
	}

	requestToDBDurationHistogram.
		WithLabelValues(category, errData).
		Observe(float64(duration.Seconds()))
}

// StoreRepositorySize сохраняет размер репозитория для указанного объекта.
func StoreRepositorySize(objectsName string, size float64) {
	repositorySizeGauge.WithLabelValues(objectsName).Set(size)
}
