package main

import (
	"time"

	"github.com/koykov/cbytecache"
)

type RequestInit struct {
	Buckets          uint                  `json:"buckets"`
	Capacity         cbytecache.MemorySize `json:"capacity"`
	ExpireIntervalNS time.Duration         `json:"expire_interval_ns"`
	EvictIntervalNS  time.Duration         `json:"evict_interval_ns"`
	VacuumIntervalNS time.Duration         `json:"vacuum_interval_ns"`
	DumpIntervalNS   time.Duration         `json:"dump_interval_ns"`
	ExpireInterval   string                `json:"expire_interval"`
	EvictInterval    string                `json:"evict_interval"`
	VacuumInterval   string                `json:"vacuum_interval"`
	DumpInterval     string                `json:"dump_interval"`
	VacuumRatio      float64               `json:"vacuum_ratio"`
	CollisionCheck   bool                  `json:"collision_check"`

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
	if conf.ExpireInterval = r.ExpireIntervalNS; conf.ExpireInterval == 0 {
		conf.ExpireInterval, _ = time.ParseDuration(r.ExpireInterval)
	}
	if conf.EvictInterval = r.EvictIntervalNS; conf.EvictInterval == 0 {
		conf.EvictInterval, _ = time.ParseDuration(r.EvictInterval)
	}
	if conf.VacuumInterval = r.VacuumIntervalNS; conf.VacuumInterval == 0 {
		conf.VacuumInterval, _ = time.ParseDuration(r.VacuumInterval)
	}
	if conf.DumpInterval = r.DumpIntervalNS; conf.DumpInterval == 0 {
		conf.DumpInterval, _ = time.ParseDuration(r.DumpInterval)
	}
	conf.VacuumRatio = r.VacuumRatio
	conf.CollisionCheck = r.CollisionCheck
	conf.Capacity = r.Capacity
}
