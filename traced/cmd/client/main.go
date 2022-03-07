package main

import (
	"flag"
	"log"
)

var (
	port = flag.Uint("port", 0, "Application port.")
)

func init() {
	flag.Parse()
	if *port == 0 {
		log.Fatalln("empty app port provided")
	}
}

func main() {}
