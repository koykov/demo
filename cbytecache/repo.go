package main

import (
	"time"

	"github.com/koykov/cbytecache"
)

type RequestInit struct {
	Buckets        uint                  `json:"buckets"`
	ExpireInterval time.Duration         `json:"expire_interval_ns"`
	EvictInterval  time.Duration         `json:"evict_interval_ns"`
	VacuumInterval time.Duration         `json:"vacuum_interval_ns"`
	VacuumRatio    float64               `json:"vacuum_ratio"`
	CollisionCheck bool                  `json:"collision_check"`
	Capacity       cbytecache.MemorySize `json:"capacity"`

	KRP uint32 `json:"krp"` // KRP - keys rotate percent

	WritersMin  uint32 `json:"writers_min"`
	WritersMax  uint32 `json:"writers_max"`
	WriterDelay uint32 `json:"writer_delay"`

	ReadersMin  uint32 `json:"readers_min"`
	ReadersMax  uint32 `json:"readers_max"`
	ReaderDelay uint32 `json:"reader_delay"`
}

func (r *RequestInit) MapConfig(conf *cbytecache.Config) {
	conf.Buckets = r.Buckets
	conf.ExpireInterval = r.ExpireInterval
	conf.EvictInterval = r.EvictInterval
	conf.VacuumInterval = r.VacuumInterval
	conf.VacuumRatio = r.VacuumRatio
	conf.CollisionCheck = r.CollisionCheck
	conf.Capacity = r.Capacity
}
