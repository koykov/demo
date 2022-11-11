package main

import (
	"math"
	"sync/atomic"
	"time"

	"github.com/koykov/laborpool"
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
	delay  uint32
	ctl    chan signal
	status status
}

func makeProducer(idx, delay uint32) *producer {
	p := &producer{
		idx:    idx,
		delay:  delay,
		ctl:    make(chan signal, 1),
		status: statusIdle,
	}
	if p.delay == 0 {
		p.delay = 50
	}
	return p
}

func (p *producer) start() {
	p.ctl <- signalInit
}

func (p *producer) stop() {
	p.ctl <- signalStop
}

func (p *producer) produce(pool *laborpool.Pool) {
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
			if p.getStatus() == statusIdle {
				return
			}
			x := Item{4, math.MaxUint64}
			w := pool.Hire()
			if w != nil {
				_ = w.Do(func() error {
					_, _ = x.Header, x.Payload
					time.Sleep(time.Duration(p.delay) * time.Nanosecond)
					return nil
				})
				pool.Fire(w)
			}
		}
	}
}

func (p *producer) setStatus(status status) {
	atomic.StoreUint32((*uint32)(&p.status), uint32(status))
}

func (p *producer) getStatus() status {
	return status(atomic.LoadUint32((*uint32)(&p.status)))
}
