package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	HttpRequestProm = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_histogram",
		Help:    "Histogram of the http request duration.",
		Buckets: prometheus.LinearBuckets(1, 1, 10),
	}, []string{"path", "method", "status"})
)
