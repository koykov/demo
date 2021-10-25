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
	MetricsKey     string                `json:"metrics_key"`

	Writers    uint32 `json:"writers"`
	Readers    uint32 `json:"readers"`
	WriteDelay uint32 `json:"write_delay"`
	ReadDelay  uint32 `json:"read_delay"`
}

func (r *RequestInit) MapConfig(conf *cbytecache.Config) {
	conf.Buckets = r.Buckets
	conf.Expire = r.Expire
	conf.Vacuum = r.Vacuum
	conf.CollisionCheck = r.CollisionCheck
	conf.MaxSize = r.MaxSize
}
