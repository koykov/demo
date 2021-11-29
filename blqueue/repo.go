package main

import (
	"time"

	"github.com/koykov/blqueue"
)

type RequestInit struct {
	Size      uint64        `json:"size"`
	Workers   uint32        `json:"workers,omitempty"`
	Heartbeat time.Duration `json:"heartbeat,omitempty"`

	WorkersMin   uint32  `json:"workers_min"`
	WorkersMax   uint32  `json:"workers_max"`
	WorkerDelay  uint32  `json:"worker_delay,omitempty"`
	WakeupFactor float32 `json:"wakeup_factor,omitempty"`
	SleepFactor  float32 `json:"sleep_factor,omitempty"`

	ProducersMin  uint32 `json:"producers_min"`
	ProducersMax  uint32 `json:"producers_max"`
	ProducerDelay uint32 `json:"producer_delay,omitempty"`

	AllowLeak bool `json:"allow_leak,omitempty"`

	Schedule []struct {
		Range        string  `json:"range,omitempty"`
		RelRange     string  `json:"rel_range,omitempty"`
		WorkersMin   uint32  `json:"workers_min,omitempty"`
		WorkersMax   uint32  `json:"workers_max,omitempty"`
		WakeupFactor float32 `json:"wakeup_factor,omitempty"`
		SleepFactor  float32 `json:"sleep_factor,omitempty"`
	} `json:"schedule,omitempty"`
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
