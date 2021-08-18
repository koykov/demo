package main

import (
	"time"

	"github.com/koykov/cbytecache"
)

type RequestInit struct {
	Buckets        uint                  `json:"buckets"`
	Expire         time.Duration         `json:"expire_ns"`
	Vacuum         time.Duration         `json:"vacuum_ns"`
	ForceSet       bool                  `json:"force_set"`
	CollisionCheck bool                  `json:"collision_check"`
	MaxSize        cbytecache.MemorySize `json:"max_size"`
	MetricsKey     string                `json:"metrics_key"`

	Writers uint32 `json:"writers"`
	Readers uint32 `json:"readers"`
}

func (r *RequestInit) MapConfig(conf *cbytecache.Config) {
	conf.Buckets = r.Buckets
	conf.Expire = r.Expire
	conf.Vacuum = r.Vacuum
	conf.ForceSet = r.ForceSet
	conf.CollisionCheck = r.CollisionCheck
	conf.MaxSize = r.MaxSize
}
