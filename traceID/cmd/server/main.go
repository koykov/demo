package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/koykov/traceID"
	"github.com/koykov/traceID/broadcaster"
)

var (
	port  = flag.Uint("port", 0, "Application port.")
	tport = flag.Uint("tport", 0, "Trace daemon port.")
	cport = flag.String("cport", "", "Client applications port separated by comma.")
)

func init() {
	flag.Parse()
	if *port == 0 {
		log.Fatalln("empty app port provided")
	}
	if *tport == 0 {
		log.Fatalln("empty traced port provided")
	}
	if len(*cport) == 0 {
		log.Fatalln("empty client applications ports provided")
	}
	cports := strings.Split(*cport, ",")
	for i := 0; i < len(cports); i++ {
		cpool = append(cpool, fmt.Sprintf("http://:%s", cports[i]))
	}
	rand.Seed(time.Now().UnixNano())

	bc := broadcaster.HTTP{Addr: fmt.Sprintf("http://:%d/post-msg", *tport)}
	traceID.RegisterBroadcaster(&bc)
}

func main() {
	addr := fmt.Sprintf(":%d", *port)
	h := ServerHTTP{}
	log.Printf("starting HTTP server at '%s'\n", addr)
	if err := http.ListenAndServe(addr, &h); err != nil {
		log.Fatalf("couldn't start HTTP server: '%s'\n", err.Error())
	}
}
