package main

import (
	"flag"
	"log"
)

var (
	port  = flag.Uint("port", 0, "Application port.")
	caddr = flag.String("caddr", "[]", "Client applications list separated by comma.")
)

func init() {
	flag.Parse()
	if *port == 0 {
		log.Fatalln("empty app port provided")
	}
}

func main() {}
