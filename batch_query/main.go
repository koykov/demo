package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	bqh    *BQHTTP
	hport  = flag.Int("hport", 8080, "HTTP port")
	pport  = flag.Int("pport", 8081, "Prometheus port")
	pfport = flag.Int("pfport", 8082, "pprof port")

	krepo keysRepo
)

func init() {
	flag.Parse()
	bqh = NewBQHTTP()
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
	if err := http.ListenAndServe(haddr, bqh); err != nil {
		log.Fatal(err)
	}
}
