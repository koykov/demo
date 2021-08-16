package main

import "github.com/koykov/cbytecache"

type demoCache struct {
	key   string
	cache *cbytecache.CByteCache

	writers,
	readers uint32
}

func (d *demoCache) Run() {
	// todo implement me
}

func (d *demoCache) String() string {
	// todo implement me
	return ""
}
