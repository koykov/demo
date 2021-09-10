package main

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/koykov/cbytecache"
)

const (
	statusIdle   status = 0
	statusActive status = 1

	signalInit signal = 0
	signalStop signal = 1

	maxIndex = 1e9
)

type status uint32
type signal uint32

type demoCache struct {
	key    string
	config *cbytecache.Config
	cache  *cbytecache.CByteCache

	offPtr uint64

	writers,
	readers uint32

	writersPool []*writer
	readersPool []*reader

	cancelFnExpire  context.CancelFunc
	cancelFnCounter context.CancelFunc
}

func (d *demoCache) Run() {
	d.offPtr = 1

	d.writersPool = make([]*writer, d.writers)
	for i := 0; i < int(d.writers); i++ {
		d.writersPool[i] = makeWriter(uint32(i), &d.offPtr)
	}
	for i := 0; i < int(d.writers); i++ {
		go d.writersPool[i].run(d.cache)
		d.writersPool[i].start()
	}

	d.readersPool = make([]*reader, d.readers)
	for i := 0; i < int(d.readers); i++ {
		d.readersPool[i] = makeReader(uint32(i), &d.offPtr)
	}
	for i := 0; i < int(d.readers); i++ {
		go d.readersPool[i].run(d.cache)
		d.readersPool[i].start()
	}

	var (
		clockExpire, clockCounter context.Context
	)
	clockExpire, d.cancelFnExpire = context.WithCancel(context.Background())
	tickerExpire := time.NewTicker(d.config.Expire)
	go func(ctx context.Context) {
		for {
			select {
			case <-tickerExpire.C:
				atomic.StoreUint64(&d.offPtr, 1)
			case <-ctx.Done():
				return
			}
		}
	}(clockExpire)

	clockCounter, d.cancelFnCounter = context.WithCancel(context.Background())
	tickerCounter := time.NewTicker(d.config.Expire)
	go func(ctx context.Context) {
		for {
			select {
			case <-tickerCounter.C:
				atomic.AddUint64(&d.offPtr, 1)
			case <-ctx.Done():
				return
			}
			time.Sleep(time.Microsecond * 1)
		}
	}(clockCounter)
}

func (d *demoCache) Stop() {
	d.cancelFnExpire()
	d.cancelFnCounter()
	// todo implement me
}

func (d *demoCache) String() string {
	// todo implement me
	return ""
}
