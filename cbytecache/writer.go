package main

import (
	"math/rand"
	"sync/atomic"

	"github.com/koykov/bytebuf"
	"github.com/koykov/cbytecache"
)

type writer struct {
	idx    uint32
	status status
	ctl    chan signal
	buf    bytebuf.ChainBuf
	offPtr *uint64
}

func makeWriter(idx uint32, offPtr *uint64) *writer {
	w := &writer{
		idx:    idx,
		status: statusIdle,
		ctl:    make(chan signal, 1),
		offPtr: offPtr,
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
			i := rand.Intn(int(atomic.LoadUint64(w.offPtr)))
			w.buf.Reset().WriteStr("key").WriteInt(int64(i))
			_ = cache.Set(w.buf.String(), getTestBody(i))
		}
	}
}

func (w *writer) setStatus(status status) {
	atomic.StoreUint32((*uint32)(&w.status), uint32(status))
}

func (w *writer) getStatus() status {
	return status(atomic.LoadUint32((*uint32)(&w.status)))
}
