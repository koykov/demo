package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/koykov/traceID"
	"github.com/koykov/traceID/broadcaster"
)

var (
	port   = flag.Uint("port", 0, "Application port.")
	cbport = flag.Uint("cbport", 0, "Callback application port.")
	pbport = flag.Uint("pbport", 0, "Postback application port.")
	tport  = flag.Uint("tport", 0, "Trace daemon port.")
	cport  = flag.String("cport", "", "Client applications port separated by comma.")

	logger = log.New(os.Stdout, "", log.LstdFlags)

	i10n chan os.Signal
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

	i10n = make(chan os.Signal, 1)
	signal.Notify(i10n, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
}

func main() {
	go func() {
		addr := fmt.Sprintf(":%d", *port)
		h := ServerHTTP{
			PortCB: *cbport,
			PortPB: *pbport,
		}
		log.Printf("starting HTTP server at '%s'\n", addr)
		if err := http.ListenAndServe(addr, &h); err != nil {
			log.Fatalf("couldn't start HTTP server: '%s'\n", err.Error())
		}
	}()

	go func() {
		addr := fmt.Sprintf(":%d", *cbport)
		h := CallbackHTTP{}
		log.Printf("starting Callback server at '%s'\n", addr)
		if err := http.ListenAndServe(addr, &h); err != nil {
			log.Fatalf("couldn't start Callback server: '%s'\n", err.Error())
		}
	}()

	go func() {
		addr := fmt.Sprintf(":%d", *pbport)
		h := PostbackHTTP{}
		log.Printf("starting Postback server at '%s'\n", addr)
		if err := http.ListenAndServe(addr, &h); err != nil {
			log.Fatalf("couldn't start Postback server: '%s'\n", err.Error())
		}
	}()

	<-i10n
}
