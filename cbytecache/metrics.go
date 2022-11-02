package main

import "github.com/prometheus/client_golang/prometheus"

var (
	writerIdle, writerActive, readerIdle, readerActive *prometheus.GaugeVec
)

func init() {
	writerIdle = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "cbytecache_writers_idle",
		Help: "Indicates how many writers idle.",
	}, []string{"cache"})
	writerActive = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "cbytecache_writers_active",
		Help: "Indicates how many writers active.",
	}, []string{"cache"})

	readerIdle = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "cbytecache_readers_idle",
		Help: "Indicates how many readers idle.",
	}, []string{"cache"})
	readerActive = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "cbytecache_readers_active",
		Help: "Indicates how many readers active.",
	}, []string{"cache"})

	prometheus.MustRegister(writerIdle, writerActive, readerIdle, readerActive)
}

func WritersInitMetric(cache string, up, idle uint32) {
	writerActive.WithLabelValues(cache).Add(float64(up))
	writerIdle.WithLabelValues(cache).Add(float64(idle))
}

func WriterStartMetric(cache string) {
	writerActive.WithLabelValues(cache).Inc()
	writerIdle.WithLabelValues(cache).Add(-1)
}

func WriterStopMetric(cache string) {
	writerIdle.WithLabelValues(cache).Inc()
	writerActive.WithLabelValues(cache).Add(-1)
}

func ReadersInitMetric(cache string, up, idle uint32) {
	readerActive.WithLabelValues(cache).Add(float64(up))
	readerIdle.WithLabelValues(cache).Add(float64(idle))
}

func ReaderStartMetric(cache string) {
	readerActive.WithLabelValues(cache).Inc()
	readerIdle.WithLabelValues(cache).Add(-1)
}

func ReaderStopMetric(cache string) {
	readerIdle.WithLabelValues(cache).Inc()
	readerActive.WithLabelValues(cache).Add(-1)
}
