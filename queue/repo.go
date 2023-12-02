package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/koykov/clock"
	q "github.com/koykov/queue"
)

type RequestInit struct {
	Capacity  uint64        `json:"capacity"`
	Instances uint32        `json:"instances"`
	Workers   uint32        `json:"workers,omitempty"`
	Heartbeat time.Duration `json:"heartbeat,omitempty"`

	WorkersMin      uint32  `json:"workers_min"`
	WorkersMax      uint32  `json:"workers_max"`
	WorkerDelay     uint32  `json:"worker_delay,omitempty"`
	WakeupFactor    float32 `json:"wakeup_factor,omitempty"`
	SleepFactor     float32 `json:"sleep_factor,omitempty"`
	SleepThreshold  uint32  `json:"sleep_threshold"`
	WorkersSchedule []struct {
		Range        string  `json:"range,omitempty"`
		RelRange     string  `json:"rel_range,omitempty"`
		WorkersMin   uint32  `json:"workers_min,omitempty"`
		WorkersMax   uint32  `json:"workers_max,omitempty"`
		WakeupFactor float32 `json:"wakeup_factor,omitempty"`
		SleepFactor  float32 `json:"sleep_factor,omitempty"`
	} `json:"workers_schedule,omitempty"`

	ProducersMin      uint32 `json:"producers_min"`
	ProducersMax      uint32 `json:"producers_max"`
	ProducerDelay     uint32 `json:"producer_delay,omitempty"`
	ProducersSchedule []struct {
		Range     string `json:"range,omitempty"`
		RelRange  string `json:"rel_range,omitempty"`
		Producers uint32 `json:"producers,omitempty"`
	} `json:"producers_schedule,omitempty"`

	AllowLeak         bool   `json:"allow_leak,omitempty"`
	LeakDirection     string `json:"leak_direction"`
	AllowDeadline     bool   `json:"allow_deadline,omitempty"`
	FrontLeakAttempts uint32 `json:"front_leak_attempts"`
	Dump              *struct {
		Capacity uint64 `json:"capacity"`
		Flush    int64  `json:"flush"`
		Buffer   uint64 `json:"buffer,omitempty"`
	} `json:"dump,omitempty"`
	Restore *struct {
		Check     int64   `json:"check"`
		Postpone  int64   `json:"postpone"`
		AllowRate float32 `json:"allow_rate"`
	} `json:"restore,omitempty"`
	DelayNs uint64 `json:"delay_ns,omitempty"`

	QoS *struct {
		Algo   string `json:"algo"`
		Egress struct {
			Capacity      uint64 `json:"capacity"`
			Instances     uint32 `json:"instances"`
			Workers       uint32 `json:"workers"`
			IdleThreshold uint32 `json:"idle_threshold"`
			IdleTimeout   int64  `json:"idle_timeout"`
		} `json:"egress"`
		Queues []struct {
			Name          string `json:"name,omitempty"`
			Capacity      uint64 `json:"capacity"`
			Weight        uint64 `json:"weight"`
			IngressWeight uint64 `json:"ingress_weight"`
			EgressWeight  uint64 `json:"egress_weight"`
		} `json:"queues"`
	} `json:"qos,omitempty"`
}

func (r *RequestInit) MapConfig(conf *q.Config) {
	conf.Capacity = r.Capacity
	conf.Instances = r.Instances
	conf.Workers = r.Workers
	conf.HeartbeatInterval = r.Heartbeat
	conf.WorkersMin = r.WorkersMin
	conf.WorkersMax = r.WorkersMax
	conf.WakeupFactor = r.WakeupFactor
	conf.SleepFactor = r.SleepFactor
	conf.SleepThreshold = r.SleepThreshold
	conf.DelayInterval = time.Duration(r.DelayNs)
	if r.LeakDirection == "front" {
		conf.LeakDirection = q.LeakDirectionFront
		conf.FrontLeakAttempts = r.FrontLeakAttempts
	}
	if len(r.WorkersSchedule) > 0 {
		now := time.Now()
		s := q.NewSchedule()
		for _, rule := range r.WorkersSchedule {
			var r1 string
			if r1 = rule.Range; len(r1) == 0 {
				if p := strings.Split(rule.RelRange, "-"); len(p) == 2 {
					var (
						d0, d1 time.Duration
						err    error
					)
					if d0, err = clock.Relative(p[0]); err != nil {
						fmt.Println("bad range", rule.RelRange, "err", err)
						continue
					}
					if d1, err = clock.Relative(p[1]); err != nil {
						fmt.Println("bad range", rule.RelRange, "err", err)
						continue
					}
					now0, now1 := now.Add(d0), now.Add(d1)
					r0 := fmt.Sprintf("%02d:%02d:%02d-%02d:%02d:%02d", now0.Hour(), now0.Minute(), now0.Second(),
						now1.Hour(), now1.Minute(), now1.Second())
					params := q.ScheduleParams{
						WorkersMin:   rule.WorkersMin,
						WorkersMax:   rule.WorkersMax,
						WakeupFactor: rule.WakeupFactor,
						SleepFactor:  rule.SleepFactor,
					}
					if err = s.AddRange(r0, params); err != nil {
						fmt.Println("error", err, "caught on adding range", r0)
						continue
					}
				}
			}
		}
		conf.Schedule = s
	}
}

func (r *RequestInit) MapInternalQueue(queue *demoQueue) {
	if len(r.ProducersSchedule) > 0 {
		now := time.Now()
		s := q.NewSchedule()
		for _, rule := range r.ProducersSchedule {
			var r1 string
			if r1 = rule.Range; len(r1) == 0 {
				if p := strings.Split(rule.RelRange, "-"); len(p) == 2 {
					var (
						d0, d1 time.Duration
						err    error
					)
					if d0, err = clock.Relative(p[0]); err != nil {
						fmt.Println("bad range", rule.RelRange, "err", err)
						continue
					}
					if d1, err = clock.Relative(p[1]); err != nil {
						fmt.Println("bad range", rule.RelRange, "err", err)
						continue
					}
					now0, now1 := now.Add(d0), now.Add(d1)
					r0 := fmt.Sprintf("%02d:%02d:%02d-%02d:%02d:%02d", now0.Hour(), now0.Minute(), now0.Second(),
						now1.Hour(), now1.Minute(), now1.Second())
					params := q.ScheduleParams{
						WorkersMin: rule.Producers,
						WorkersMax: rule.Producers + 1,
					}
					if params.WorkersMin > r.ProducersMax {
						params.WorkersMin = r.ProducersMax
						params.WorkersMax = r.ProducersMax + 1
					}
					if err = s.AddRange(r0, params); err != nil {
						fmt.Println("error", err, "caught on adding range", r0)
						continue
					}
				}
			}
		}
		queue.schedule = s
	}
}
