package main

import (
	"log"
	"math/rand"
	"sort"
	"sync"
	"time"

	"github.com/koykov/byteconv"
)

type ckey struct {
	key    string
	expire uint32
}

type keyRegistry struct {
	mux  sync.RWMutex
	keys []ckey
}

var (
	chars = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	keys  keyRegistry
)

func (r *keyRegistry) get(newPercent int) string {
	r.mux.RLock()
	defer r.mux.RUnlock()
	if rand.Intn(100) < newPercent {
		l := rand.Intn(16) + 16
		b := make([]byte, l)
		for i := 0; i < l; i++ {
			b[i] = chars[rand.Intn(len(chars)-1)]
		}
		return byteconv.B2S(b)
	} else {
		if len(r.keys) <= 1 {
			return ""
		}
		i := rand.Intn(len(r.keys) - 1)
		key := &r.keys[i]
		return key.key
	}
}

func (r *keyRegistry) set(key string, expire time.Duration) {
	r.mux.Lock()
	r.keys = append(r.keys, ckey{
		key:    key,
		expire: uint32(time.Now().Unix() + int64(expire.Seconds())),
	})
	r.mux.Unlock()
}

func (r *keyRegistry) bulkEvict() {
	r.mux.Lock()
	defer r.mux.Unlock()
	now := uint32(time.Now().Unix())
	l := len(r.keys)
	if l == 0 {
		return
	}
	_ = r.keys[l-1]
	z := sort.Search(l-1, func(i int) bool {
		return now <= r.keys[i].expire
	})
	if z == 0 {
		return
	}
	log.Printf("found %d expired keys", z)
	copy(r.keys[0:], r.keys[z:])
	r.keys = r.keys[l-z:]
}

// func (r *keyRegistry) Listen(entry ttlcache.Entry) error {
// 	key := entry.Key
// 	_ = key // todo remove key from registry
// 	return nil
// }

func (r *keyRegistry) Close() error { return nil }
