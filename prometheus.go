package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func AddPrometheusMetrics() {
	prometheus.MustRegister(ResponseTimeHistogram)
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}

var (
	UserRegisteredCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "number_of_registered_users",
		Help: "The total number of registered users",
	})
	CakesGiven = promauto.NewCounter(prometheus.CounterOpts{
		Name: "number_of_cakes_given",
		Help: "Number of cakes given",
	})
)

var buckets = []float64{.0000001, .00001, .0001, .001, .005, .01, .05, .1, .2, .25, .5, 1, 2.5, 5, 10}

var ResponseTimeHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Namespace: "namespace",
	Name:      "http_server_request_duration_seconds",
	Help:      "Histogram of response time for handler in miliseconds",
	Buckets:   buckets,
}, []string{"route", "method", "status_code"})

