package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTPRequestDuration measures the duration of HTTP requests by method and path.
	HTTPRequestDuration *prometheus.HistogramVec

	// ServiceMethodDuration measures the duration of service method executions by method name.
	ServiceMethodDuration *prometheus.HistogramVec

	// RepoMethodDuration measures the duration of repository method executions by method name.
	RepoMethodDuration *prometheus.HistogramVec

	// EnricherRequestDuration measures the duration of external enricher API requests by type.
	EnricherRequestDuration *prometheus.HistogramVec
)

// InitMetrics инициализирует метрики Prometheus.
func InitMetrics() {
	RepoMethodDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "repo_method_duration_seconds",
		Help:    "Duration of repository method executions by method name.",
		Buckets: prometheus.DefBuckets,
	}, []string{"method"})

	EnricherRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "enricher_request_duration_seconds",
		Help:    "Duration of external enricher API requests by type.",
		Buckets: prometheus.DefBuckets,
	}, []string{"type"})

	HTTPRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "Duration of HTTP requests by method and path.",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "path"})

	ServiceMethodDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "service_method_duration_seconds",
		Help:    "Duration of service method executions by method name.",
		Buckets: prometheus.DefBuckets,
	}, []string{"method"})
}
