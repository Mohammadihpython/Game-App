package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"sync"
)

var (
	once                 sync.Once
	requestCounter       *prometheus.CounterVec
	requestStatusCounter *prometheus.CounterVec
	requestDuration      *prometheus.HistogramVec
)

func Init() {
	once.Do(func() {
		requestStatusCounter = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "request_status_counter",
			},
			[]string{"code", "method", "endpoint"},
		)

		requestCounter = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "login_request_counter",
				Help: "Number of login requests",
			},
			[]string{"service", "endpoint"},
		)
		requestDuration = prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "Duration of http requests",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"service", "endpoint"},
		)
		prometheus.MustRegister(requestCounter, requestDuration, requestStatusCounter)

	})

}

// Expose metrics for use in other packages

func RequestCounter() *prometheus.CounterVec {
	return requestCounter
}

func RequestStatusCounter() *prometheus.CounterVec {
	return requestStatusCounter
}

func RequestDuration() *prometheus.HistogramVec {
	return requestDuration
}
