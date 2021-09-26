package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func AddPrometheusMetrics() {
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
