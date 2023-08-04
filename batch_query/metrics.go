package main

import "github.com/prometheus/client_golang/prometheus"

var producerIdle, producerActive *prometheus.GaugeVec

func init() {
	producerIdle = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "batch_query_producers_idle",
		Help: "Indicates how many producers idle.",
	}, []string{"query"})
	producerActive = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "batch_query_producers_active",
		Help: "Indicates how many producers active.",
	}, []string{"query"})

	prometheus.MustRegister(producerIdle, producerActive)
}

func ProducersInitMetric(query string, up, idle uint32) {
	producerActive.WithLabelValues(query).Add(float64(up))
	producerIdle.WithLabelValues(query).Add(float64(idle))
}

func ProducerStartMetric(query string) {
	producerActive.WithLabelValues(query).Inc()
	producerIdle.WithLabelValues(query).Add(-1)
}

func ProducerStopMetric(query string) {
	producerIdle.WithLabelValues(query).Inc()
	producerActive.WithLabelValues(query).Add(-1)
}
