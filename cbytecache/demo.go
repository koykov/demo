package main

import (
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
	cache  *cbytecache.CByteCache

	writers,
	readers uint32

	writersPool []*writer
	readersPool []*reader
}

func (d *demoCache) Run() {
	d.writersPool = make([]*writer, d.writers)
	for i := 0; i < int(d.writers); i++ {
		d.writersPool[i] = makeWriter(uint32(i), d.config)
	}
	for i := 0; i < int(d.writers); i++ {
		go d.writersPool[i].run(d.cache)
		d.writersPool[i].start()
	}

	d.readersPool = make([]*reader, d.readers)
	for i := 0; i < int(d.readers); i++ {
		d.readersPool[i] = makeReader(uint32(i), d.config)
	}
	for i := 0; i < int(d.readers); i++ {
		go d.readersPool[i].run(d.cache)
		d.readersPool[i].start()
	}
}

func (d *demoCache) Stop() {
	// todo implement me
}

func (d *demoCache) String() string {
	// todo implement me
	return ""
}
