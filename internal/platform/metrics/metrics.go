package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var HTTPRequestsTotal = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests processed",
	},
	[]string{"method", "path", "status"},
)

var HTTPRequestDuration = promauto.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "Histogram of response latency for HTTP requests",
		Buckets: []float64{0.001, 0.005, 0.015, 0.05, 0.1, 0.25, 0.5, 1.0, 2.5}, // 1ms up to 2.5s
	},
	[]string{"method", "path"},
)

var ActiveSeatHolds = promauto.NewGauge(
	prometheus.GaugeOpts{
		Name: "active_seat_holds",
		Help: "Current number of active seat holds in Redis",
	},
)
