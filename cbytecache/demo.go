package main

import (
	"github.com/koykov/cbytecache"
)

type status uint32
type signal uint32

const (
	statusIdle   status = 0
	statusActive status = 1

	signalInit signal = 0
	signalStop signal = 1
)

type demoCache struct {
	key   string
	cache *cbytecache.CByteCache

	writers,
	readers uint32

	writersPool []*writer
}

func (d *demoCache) Run() {
	d.writersPool = make([]*writer, d.writers)
	for i := 0; i < int(d.writers); i++ {
		d.writersPool[i] = makeWriter(uint32(i))
	}
	for i := 0; i < int(d.writers); i++ {
		go d.writersPool[i].run(d.cache)
		d.writersPool[i].start()
	}
}

func (d *demoCache) Stop() {
	// todo implement me
}

func (d *demoCache) String() string {
	// todo implement me
	return ""
}
