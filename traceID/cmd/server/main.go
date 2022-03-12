package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

var (
	port  = flag.Uint("port", 0, "Application port.")
	cport = flag.String("cport", "", "Client applications port separated by comma.")
)

func init() {
	flag.Parse()
	if *port == 0 {
		log.Fatalln("empty app port provided")
	}
	if len(*cport) == 0 {
		log.Fatalln("empty client applications ports provided")
	}
	cports := strings.Split(*cport, ",")
	for i := 0; i < len(cports); i++ {
		cpool = append(cpool, fmt.Sprintf("http://:%s", cports[i]))
	}
	rand.Seed(time.Now().UnixNano())
}

func main() {
	addr := fmt.Sprintf(":%d", *port)
	h := ServerHTTP{}
	log.Printf("starting HTTP server at '%s'\n", addr)
	if err := http.ListenAndServe(addr, &h); err != nil {
		log.Fatalf("couldn't start HTTP server: '%s'\n", err.Error())
	}
}
