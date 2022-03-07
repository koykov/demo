package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var (
	port = flag.Uint("port", 0, "Application port.")
)

func init() {
	flag.Parse()
	if *port == 0 {
		log.Fatalln("empty app port provided")
	}
	rand.Seed(time.Now().UnixNano())
}

func main() {
	addr := fmt.Sprintf(":%d", *port)
	h := ClientHTTP{}
	log.Printf("starting HTTP client at '%s'\n", addr)
	if err := http.ListenAndServe(addr, &h); err != nil {
		log.Fatalf("couldn't start HTTP client: '%s'\n", err.Error())
	}
}
