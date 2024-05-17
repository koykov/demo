package main

import (
	"time"

	"github.com/koykov/ttlcache"
)

type RequestInit struct {
	Buckets         uint          `json:"buckets"`
	Size            uint64        `json:"size"`
	TTLIntervalNS   time.Duration `json:"ttl_interval_ns"`
	EvictIntervalNS time.Duration `json:"evict_interval_ns"`
	EvictWorkers    uint          `json:"evict_workers"`

	DeletePercent uint `json:"delete_percent"`

	WritersMin  uint32 `json:"writers_min"`
	WritersMax  uint32 `json:"writers_max"`
	WriterDelay uint32 `json:"writer_delay"`

	ReadersMin  uint32 `json:"readers_min"`
	ReadersMax  uint32 `json:"readers_max"`
	ReaderDelay uint32 `json:"reader_delay"`
}

func (r *RequestInit) MapConfig(conf *ttlcache.Config[entry]) {
	conf.Buckets = r.Buckets
	conf.Size = r.Size
	conf.TTLInterval = r.TTLIntervalNS
	conf.EvictInterval = r.EvictIntervalNS
	conf.EvictWorkers = r.EvictWorkers
	if r.DeletePercent > 100 {
		r.DeletePercent = 0
	}
}
