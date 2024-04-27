package main

import (
	"bytes"
	"log"
	"math/rand"
	"sync/atomic"
	"time"

	"github.com/koykov/ttlcache"
)

type reader struct {
	idx    uint32
	status status
	config *ttlcache.Config[entry]
	req    *RequestInit
	ctl    chan signal
	dst    entry
}

func makeReader(idx uint32, config *ttlcache.Config[entry], req *RequestInit) *reader {
	r := &reader{
		idx:    idx,
		status: statusIdle,
		config: config,
		req:    req,
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

func (r *reader) run(cache *ttlcache.Cache[entry]) {
	var err error
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
			if key := keys.get(0); len(key) > 0 {
				if r.dst, err = cache.Get(key); err == nil {
					e := testData[r.dst.c]
					if !bytes.Equal(r.dst.b, e) {
						log.Println("bad answer")
					}
					if r.req.DeletePercent > 0 && uint(rand.Intn(100)) < r.req.DeletePercent {
						_ = cache.Delete(key)
					}
				} else {
					// log.Println("err", err, "len", len(r.dst))
				}
				if delay := r.req.ReaderDelay; delay > 0 {
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
