package main

import (
	"time"

	"github.com/koykov/blqueue"
)

type RequestInit struct {
	Size      uint64        `json:"size"`
	Workers   uint32        `json:"workers"`
	Heartbeat time.Duration `json:"heartbeat"`

	WorkersMin   uint32  `json:"workers_min"`
	WorkersMax   uint32  `json:"workers_max"`
	WorkerDelay  uint32  `json:"worker_delay"`
	WakeupFactor float32 `json:"wakeup_factor"`
	SleepFactor  float32 `json:"sleep_factor"`

	MetricsKey string `json:"metrics_key"`

	ProducersMin  uint32 `json:"producers_min"`
	ProducersMax  uint32 `json:"producers_max"`
	ProducerDelay uint32 `json:"producer_delay"`

	AllowLeak bool `json:"allow_leak"`
}

func (r *RequestInit) MapConfig(conf *blqueue.Config) {
	conf.Size = r.Size
	conf.Workers = r.Workers
	conf.Heartbeat = r.Heartbeat
	conf.WorkersMin = r.WorkersMin
	conf.WorkersMax = r.WorkersMax
	conf.WakeupFactor = r.WakeupFactor
	conf.SleepFactor = r.SleepFactor
}
