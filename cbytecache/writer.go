package main

import (
	"sync/atomic"
	"time"

	"github.com/koykov/cbytecache"
)

type writer struct {
	idx    uint32
	status status
	config *cbytecache.Config
	rawReq *RequestInit
	ctl    chan signal
}

func makeWriter(idx uint32, config *cbytecache.Config, rawReq *RequestInit) *writer {
	w := &writer{
		idx:    idx,
		status: statusIdle,
		config: config,
		rawReq: rawReq,
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

func (w *writer) run(cache *cbytecache.CByteCache) {
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
			key := keys.get(int(w.rawReq.WriterKRP))
			_ = cache.Set(key, getTestBody())
			keys.set(key, w.config.Expire)
			if delay := w.rawReq.WriteDelay; delay > 0 {
				time.Sleep(time.Duration(delay))
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
