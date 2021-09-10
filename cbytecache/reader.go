package main

import (
	"math/rand"
	"sync/atomic"

	"github.com/koykov/bytealg"
	"github.com/koykov/cbytecache"
)

type reader struct {
	idx    uint32
	status status
	ctl    chan signal
	buf    bytealg.ChainBuf
	dst    []byte
	offPtr *uint64
}

func makeReader(idx uint32, offPtr *uint64) *reader {
	r := &reader{
		idx:    idx,
		status: statusIdle,
		ctl:    make(chan signal, 1),
		offPtr: offPtr,
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
			i := rand.Intn(int(atomic.LoadUint64(r.offPtr)))
			r.buf.Reset().WriteStr("key").WriteInt(int64(i))
			r.dst, _ = cache.GetTo(r.dst[:0], r.buf.String())
		}
	}
}

func (r *reader) setStatus(status status) {
	atomic.StoreUint32((*uint32)(&r.status), uint32(status))
}

func (r *reader) getStatus() status {
	return status(atomic.LoadUint32((*uint32)(&r.status)))
}
