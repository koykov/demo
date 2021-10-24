package main

import (
	"sync/atomic"

	"github.com/koykov/cbytecache"
)

type reader struct {
	idx    uint32
	status status
	config *cbytecache.Config
	ctl    chan signal
	dst    []byte
}

func makeReader(idx uint32, config *cbytecache.Config) *reader {
	r := &reader{
		idx:    idx,
		status: statusIdle,
		config: config,
		ctl:    make(chan signal, 1),
	}
	return r
}

func (r *reader) start() {
	r.ctl <- signalInit
}

func (r *reader) stop() {
	r.ctl <- signalStop
}

func (r *reader) run(cache *cbytecache.CByteCache) {
	for {
		select {
		case cmd := <-r.ctl:
			switch cmd {
			case signalInit:
				r.setStatus(statusActive)
			case signalStop:
				r.setStatus(statusIdle)
				return
			}
		default:
			if r.getStatus() == statusIdle {
				return
			}
			key := keys.get(10)
			r.dst, _ = cache.GetTo(r.dst[:0], key)
		}
	}
}

func (r *reader) setStatus(status status) {
	atomic.StoreUint32((*uint32)(&r.status), uint32(status))
}

func (r *reader) getStatus() status {
	return status(atomic.LoadUint32((*uint32)(&r.status)))
}
