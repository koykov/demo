package main

import (
	"bytes"
	"log"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/koykov/cbytecache"
	"github.com/koykov/fastconv"
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
			if key := keys.get(int(r.rawReq.ReaderKRP)); len(key) > 0 {
				if r.dst, err = cache.GetTo(r.dst[:0], key); err == nil && len(r.dst) > 0 {
					ri := r.dst[:1]
					if i, err := strconv.ParseInt(fastconv.B2S(ri), 10, 64); err == nil {
						e := testData[i]
						if !bytes.Equal(r.dst, e) {
							log.Println("bad answer")
						}
					}
				} else {
					// log.Println("err", err, "len", len(r.dst))
				}
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
