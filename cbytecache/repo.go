package main

import (
	"time"

	"github.com/koykov/cbytecache"
)

type RequestInit struct {
	Buckets        uint                  `json:"buckets"`
	Expire         time.Duration         `json:"expire_ns"`
	Vacuum         time.Duration         `json:"vacuum_ns"`
	CollisionCheck bool                  `json:"collision_check"`
	MaxSize        cbytecache.MemorySize `json:"max_size"`

	Writers     uint32 `json:"writers"`
	WriterKRP   uint32 `json:"writer_krp"` // KRP - keys rotate percent
	WriterDelay uint32 `json:"writer_delay"`
	Readers     uint32 `json:"readers"`
	ReaderKRP   uint32 `json:"reader_krp"` // KRP - keys rotate percent
	ReaderDelay uint32 `json:"reader_delay"`
}

func (r *RequestInit) MapConfig(conf *cbytecache.Config) {
	conf.Buckets = r.Buckets
	conf.Expire = r.Expire
	conf.Vacuum = r.Vacuum
	conf.CollisionCheck = r.CollisionCheck
	conf.MaxSize = r.MaxSize
}
