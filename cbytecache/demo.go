package main

import (
	"context"
	"errors"
	"time"

	"github.com/koykov/cbytecache"
)

const (
	statusIdle   status = 0
	statusActive status = 1

	signalInit signal = 0
	signalStop signal = 1
)

type status uint32
type signal uint32

type demoCache struct {
	key    string
	config *cbytecache.Config
	req    *RequestInit
	cache  *cbytecache.CByteCache

	writersUp uint32
	writers   []*writer

	readersUp uint32
	readers   []*reader

	cancel context.CancelFunc
}

func (d *demoCache) Run() {
	d.writers = make([]*writer, d.req.WritersMax)
	for i := 0; i < int(d.req.WritersMax); i++ {
		d.writers[i] = makeWriter(uint32(i), d.config, d.req)
	}
	for i := 0; i < int(d.req.WritersMin); i++ {
		go d.writers[i].run(d.cache)
		d.writers[i].start()
	}
	d.writersUp = d.req.WritersMin
	WritersInitMetric(d.key, d.writersUp, d.req.WritersMax-d.writersUp)

	d.readers = make([]*reader, d.req.ReadersMax)
	for i := 0; i < int(d.req.ReadersMax); i++ {
		d.readers[i] = makeReader(uint32(i), d.config, d.req)
	}
	for i := 0; i < int(d.req.ReadersMin); i++ {
		go d.readers[i].run(d.cache)
		d.readers[i].start()
	}
	d.readersUp = d.req.ReadersMin
	ReadersInitMetric(d.key, d.readersUp, d.req.ReadersMax-d.readersUp)

	var clockExpire context.Context
	clockExpire, d.cancel = context.WithCancel(context.Background())
	tickerExpire := time.NewTicker(d.config.ExpireInterval / 4)
	go func(ctx context.Context) {
		for {
			select {
			case <-tickerExpire.C:
				keys.bulkEvict()
			case <-ctx.Done():
				return
			}
		}
	}(clockExpire)
}

func (d *demoCache) WritersUp(delta uint32) error {
	if delta == 0 {
		delta = 1
	}
	if d.writersUp+delta >= d.req.WritersMax {
		return errors.New("maximum writers count reached")
	}
	c := d.writersUp
	for i := c; i < c+delta; i++ {
		go d.writers[i].run(d.cache)
		d.writers[i].start()
		d.writersUp++
		WriterStartMetric(d.key)
	}
	return nil
}

func (d *demoCache) WritersDown(delta uint32) error {
	if delta == 0 {
		delta = 1
	}
	if d.writersUp-delta < d.req.WritersMin {
		return errors.New("minimum writers count reached")
	}
	c := d.writersUp
	for i := c; i >= c-delta; i-- {
		if d.writers[i].getStatus() == statusActive {
			d.writers[i].stop()
			d.writersUp--
			WriterStopMetric(d.key)
		}
	}
	return nil
}

func (d *demoCache) ReadersUp(delta uint32) error {
	if delta == 0 {
		delta = 1
	}
	if d.writersUp+delta >= d.req.ReadersMax {
		return errors.New("maximum readers count reached")
	}
	c := d.readersUp
	for i := c; i < c+delta; i++ {
		go d.readers[i].run(d.cache)
		d.readers[i].start()
		d.readersUp++
		ReaderStartMetric(d.key)
	}
	return nil
}

func (d *demoCache) ReadersDown(delta uint32) error {
	if delta == 0 {
		delta = 1
	}
	if d.readersUp-delta < d.req.ReadersMin {
		return errors.New("minimum readers count reached")
	}
	c := d.readersUp
	for i := c; i >= c-delta; i-- {
		if d.readers[i].getStatus() == statusActive {
			d.readers[i].stop()
			d.readersUp--
			ReaderStopMetric(d.key)
		}
	}
	return nil
}

func (d *demoCache) Stop() {
	d.stop(false)
}

func (d *demoCache) ForceStop() {
	d.stop(true)
}

func (d *demoCache) stop(force bool) {
	c := d.writersUp
	for i := uint32(0); i < c; i++ {
		d.writers[i].stop()
		d.writersUp--
		WriterStopMetric(d.key)
	}

	c = d.readersUp
	for i := uint32(0); i < c; i++ {
		d.readers[i].stop()
		d.readersUp--
		ReaderStopMetric(d.key)
	}

	if force {
		// todo implement me, still use Close()
		_ = d.cache.Close()
	} else {
		_ = d.cache.Close()
	}
	d.cancel()
}

func (d *demoCache) String() string {
	// todo implement me
	return ""
}
