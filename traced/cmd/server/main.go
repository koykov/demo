package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"
)

var (
	port  = flag.Uint("port", 0, "Application port.")
	cport = flag.String("cport", "", "Client applications port separated by comma.")
	caddr []string
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
		caddr = append(caddr, fmt.Sprintf(":%s", cports[i]))
	}
	rand.Seed(time.Now().UnixNano())
}

func main() {}
