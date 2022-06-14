package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/koykov/traceID"
	"github.com/koykov/traceID/broadcaster"
	"github.com/koykov/tracemod/zeromq"
)

var (
	conf   Config
	logger = log.New(os.Stdout, "", log.LstdFlags)
	i10n   chan os.Signal
)

func init() {
	rand.Seed(time.Now().UnixNano())

	if err := conf.LoadFrom("config/demo.json"); err != nil {
		log.Fatalln(err)
	}
	for i := 0; i < len(conf.Clients); i++ {
		cpool = append(cpool, fmt.Sprintf("http://:%d", conf.Clients[i]))
	}

	var bc traceID.Broadcaster
	switch conf.Broadcaster.Handler {
	case "http":
		bc = &broadcaster.HTTP{}
	case "zeromq":
		bc = &zeromq.Broadcaster{}
	}
	bc.SetConfig(&conf.Broadcaster)
	traceID.RegisterBroadcaster(bc)

	i10n = make(chan os.Signal, 1)
	signal.Notify(i10n, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
}

func main() {
	for i := 0; i < len(conf.Clients); i++ {
		go func(port uint) {
			cmd := exec.CommandContext(context.Background(), "./bin/trace-client", "-port", fmt.Sprintf("%d", port))
			err := cmd.Run()
			if err != nil {
				log.Println(err)
			}
		}(conf.Clients[i])
	}

	go func() {
		addr := fmt.Sprintf(":%d", conf.Listen.AppPort)
		h := ServerHTTP{
			PortCB: conf.Listen.CbPort,
			PortPB: conf.Listen.PbPort,
		}
		log.Printf("starting HTTP server at '%s'\n", addr)
		if err := http.ListenAndServe(addr, &h); err != nil {
			log.Fatalf("couldn't start HTTP server: '%s'\n", err.Error())
		}
	}()

	go func() {
		addr := fmt.Sprintf(":%d", conf.Listen.CbPort)
		h := CallbackHTTP{}
		log.Printf("starting Callback server at '%s'\n", addr)
		if err := http.ListenAndServe(addr, &h); err != nil {
			log.Fatalf("couldn't start Callback server: '%s'\n", err.Error())
		}
	}()

	go func() {
		addr := fmt.Sprintf(":%d", conf.Listen.PbPort)
		h := PostbackHTTP{}
		log.Printf("starting Postback server at '%s'\n", addr)
		if err := http.ListenAndServe(addr, &h); err != nil {
			log.Fatalf("couldn't start Postback server: '%s'\n", err.Error())
		}
	}()

	<-i10n
}
