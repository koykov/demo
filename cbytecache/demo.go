package main

import "github.com/koykov/cbytecache"

type demoCache struct {
	key   string
	cache *cbytecache.CByteCache

	writers,
	readers uint32
}
