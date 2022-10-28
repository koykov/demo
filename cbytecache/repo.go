package main

import (
	"time"

	"github.com/koykov/cbytecache"
)

type RequestInit struct {
	Buckets        uint                  `json:"buckets"`
	ExpireInterval time.Duration         `json:"expire_interval_ns"`
	VacuumInterval time.Duration         `json:"vacuum_interval_ns"`
	CollisionCheck bool                  `json:"collision_check"`
	Capacity       cbytecache.MemorySize `json:"capacity"`

	Writers     uint32 `json:"writers"`
	WriterKRP   uint32 `json:"writer_krp"` // KRP - keys rotate percent
	WriterDelay uint32 `json:"writer_delay"`
	Readers     uint32 `json:"readers"`
	ReaderKRP   uint32 `json:"reader_krp"` // KRP - keys rotate percent
	ReaderDelay uint32 `json:"reader_delay"`
}

func (r *RequestInit) MapConfig(conf *cbytecache.Config) {
	conf.Buckets = r.Buckets
	conf.ExpireInterval = r.ExpireInterval
	conf.VacuumInterval = r.VacuumInterval
	conf.CollisionCheck = r.CollisionCheck
	conf.Capacity = r.Capacity
}
