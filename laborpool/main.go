package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	qh     *PoolHTTP
	hport  = flag.Int("hport", 8080, "HTTP port")
	pport  = flag.Int("pport", 8081, "Prometheus port")
	pfport = flag.Int("pfport", 8082, "pprof port")
)

func init() {
	flag.Parse()
	qh = NewPoolHTTP(*hport, *pport)
}

func main() {
	paddr := fmt.Sprintf(":%d", *pport)
	go func() {
		// registered metrics endpoint
		http.Handle("/metrics", promhttp.Handler())
		log.Printf("Start Prometheus server on address %s/metrics\n", paddr)
		if err := http.ListenAndServe(paddr, nil); err != nil {
			log.Fatal(err)
		}
	}()

	pfaddr := fmt.Sprintf("localhost:%d", *pfport)
	go func() {
		log.Printf("Start pprof server on address %s\n", pfaddr)
		if err := http.ListenAndServe(pfaddr, nil); err != nil {
			log.Fatal(err)
		}
	}()

	haddr := fmt.Sprintf(":%d", *hport)
	log.Println("Start HTTP server on address", haddr)
	if err := http.ListenAndServe(haddr, qh); err != nil {
		log.Fatal(err)
	}
}
