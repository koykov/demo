package main

import (
	"sync/atomic"
	"time"

	"github.com/koykov/ttlcache"
)

type writer struct {
	idx    uint32
	status status
	config *ttlcache.Config[entry]
	req    *RequestInit
	ctl    chan signal
}

func makeWriter(idx uint32, config *ttlcache.Config[entry], req *RequestInit) *writer {
	w := &writer{
		idx:    idx,
		status: statusIdle,
		config: config,
		req:    req,
		ctl:    make(chan signal, 1),
	}
	return w
}

func (w *writer) start() {
	w.ctl <- signalInit
}

func (w *writer) stop() {
	w.ctl <- signalStop
}

func (w *writer) run(cache *ttlcache.Cache[entry]) {
	for {
		select {
		case cmd := <-w.ctl:
			switch cmd {
			case signalInit:
				w.setStatus(statusActive)
			case signalStop:
				w.setStatus(statusIdle)
				return
			}
		default:
			if w.getStatus() == statusIdle {
				return
			}
			if key := keys.get(100); len(key) > 0 {
				_ = cache.Set(key, getTestBody())
				keys.set(key, w.config.EvictInterval)
				if delay := w.req.WriterDelay; delay > 0 {
					time.Sleep(time.Duration(delay))
				}
			}
		}
	}
}

func (w *writer) setStatus(status status) {
	atomic.StoreUint32((*uint32)(&w.status), uint32(status))
}

func (w *writer) getStatus() status {
	return status(atomic.LoadUint32((*uint32)(&w.status)))
}
