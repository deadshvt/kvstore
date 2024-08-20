package monitoring

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
)

type MetricsCollector struct {
	RequestCount    *prometheus.CounterVec
	RequestDuration *prometheus.HistogramVec
	Logger          zerolog.Logger
}

func NewMetricsCollector(logger zerolog.Logger) *MetricsCollector {
	return &MetricsCollector{
		RequestCount: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_request_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"path"}),
		RequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "http_request_duration_seconds",
				Help: "Duration of HTTP requests in seconds",
			},
			[]string{"path"}),
		Logger: logger,
	}
}

func (mc *MetricsCollector) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timer := prometheus.NewTimer(mc.RequestDuration.WithLabelValues(r.URL.Path))
		defer timer.ObserveDuration()
		mc.RequestCount.WithLabelValues(r.URL.Path).Inc()

		next.ServeHTTP(w, r)
	})
}

func (mc *MetricsCollector) Register() {
	prometheus.MustRegister(mc.RequestCount)
	prometheus.MustRegister(mc.RequestDuration)
}
