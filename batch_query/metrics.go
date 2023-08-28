package main

import "github.com/prometheus/client_golang/prometheus"

var producers *prometheus.GaugeVec

func init() {
	producers = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "batch_query_producers",
		Help: "Indicates how many producers idle.",
	}, []string{"query", "status"})

	prometheus.MustRegister(producers)
}

func ProducersInitMetric(query string, up, idle uint32) {
	producers.WithLabelValues(query, "active").Add(float64(up))
	producers.WithLabelValues(query, "idle").Add(float64(idle))
}

func ProducerStartMetric(query string) {
	producers.WithLabelValues(query, "active").Inc()
	producers.WithLabelValues(query, "idle").Add(-1)
}

func ProducerStopMetric(query string) {
	producers.WithLabelValues(query, "active").Add(-1)
	producers.WithLabelValues(query, "idle").Inc()
}
