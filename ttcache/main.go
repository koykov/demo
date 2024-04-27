package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	ch    *CacheHTTP
	hport = flag.Int("hport", 8080, "HTTP port")
	pport = flag.Int("pport", 8081, "Prometheus port")
)

func init() {
	rand.Seed(time.Now().UnixNano())
	flag.Parse()

	ch = NewCacheHTTP(*hport, *pport)
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

	haddr := fmt.Sprintf(":%d", *hport)
	log.Println("Start HTTP server on address", haddr)
	if err := http.ListenAndServe(haddr, ch); err != nil {
		log.Fatal(err)
	}
}
