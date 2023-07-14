package main

import (
	"sync/atomic"

	"github.com/koykov/batch_query"
)

type status uint32
type signal uint32

const (
	statusIdle   status = 0
	statusActive status = 1

	signalInit signal = 0
	signalStop signal = 1
)

type producer struct {
	idx    uint32
	ctl    chan signal
	status status
}

func makeProducer(idx uint32) *producer {
	p := &producer{
		idx:    idx,
		ctl:    make(chan signal, 1),
		status: statusIdle,
	}
	return p
}

func (p *producer) start() {
	p.ctl <- signalInit
}

func (p *producer) stop() {
	p.ctl <- signalStop
}

func (p *producer) produce(bq *batch_query.BatchQuery) {
	for {
		select {
		case cmd := <-p.ctl:
			switch cmd {
			case signalInit:
				p.setStatus(statusActive)
			case signalStop:
				p.setStatus(statusIdle)
				return
			}
		default:
			key := krepo.get()
			rec, err := bq.Find(key)
			_, _ = rec, err
		}
	}
}

func (p *producer) setStatus(status status) {
	atomic.StoreUint32((*uint32)(&p.status), uint32(status))
}

func (p *producer) getStatus() status {
	return status(atomic.LoadUint32((*uint32)(&p.status)))
}
