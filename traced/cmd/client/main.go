package main

import (
	"flag"
	"log"
	"math/rand"
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

func main() {}
