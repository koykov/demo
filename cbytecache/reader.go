package main

import (
	"sync/atomic"
	"time"

	"github.com/koykov/cbytecache"
)

type reader struct {
	idx    uint32
	status status
	config *cbytecache.Config
	rawReq *RequestInit
	ctl    chan signal
	dst    []byte
}

func makeReader(idx uint32, config *cbytecache.Config, rawReq *RequestInit) *reader {
	r := &reader{
		idx:    idx,
		status: statusIdle,
		config: config,
		rawReq: rawReq,
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
			if key := keys.get(int(r.rawReq.ReaderKRP)); len(key) > 0 {
				r.dst, _ = cache.GetTo(r.dst[:0], key)
				if delay := r.rawReq.ReaderDelay; delay > 0 {
					time.Sleep(time.Duration(delay))
				}
			}
		}
	}
}

func (r *reader) setStatus(status status) {
	atomic.StoreUint32((*uint32)(&r.status), uint32(status))
}

func (r *reader) getStatus() status {
	return status(atomic.LoadUint32((*uint32)(&r.status)))
}
