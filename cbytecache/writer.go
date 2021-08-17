package main

import (
	"sync/atomic"

	"github.com/koykov/cbytecache"
)

type writer struct {
	idx    uint32
	ctl    chan signal
	status status
}

func makeWriter(idx uint32) *writer {
	w := &writer{
		idx:    idx,
		ctl:    make(chan signal, 1),
		status: statusIdle,
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
			// todo write to cache
			_ = cache
		}
	}
}

func (w *writer) setStatus(status status) {
	atomic.StoreUint32((*uint32)(&w.status), uint32(status))
}

func (w *writer) getStatus() status {
	return status(atomic.LoadUint32((*uint32)(&w.status)))
}
