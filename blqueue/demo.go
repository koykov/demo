package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/koykov/blqueue"
)

type demoQueue struct {
	key   string
	queue *blqueue.Queue
	req   *RequestInit

	producersUp uint32
	producers   []*producer

	schedule *blqueue.Schedule
	schedID  int
	pmax     uint32

	cancel context.CancelFunc
}

type realtimeParams struct {
	min, max uint32
}

func (d *demoQueue) Run() {
	d.pmax = d.producersMaxDaily()

	d.producers = make([]*producer, d.pmax)
	for i := 0; i < int(d.pmax); i++ {
		d.producers[i] = makeProducer(uint32(i), d.req.ProducerDelay)
	}
	for i := 0; i < int(d.req.ProducersMin); i++ {
		go d.producers[i].produce(d.queue)
		d.producers[i].start()
	}
	d.producersUp = d.req.ProducersMin

	producerActive.WithLabelValues(d.key).Add(float64(d.producersUp))
	producerIdle.WithLabelValues(d.key).Add(float64(d.pmax - d.producersUp))

	d.schedID = -1

	var ctx context.Context
	ctx, d.cancel = context.WithCancel(context.Background())
	ticker := time.NewTicker(time.Millisecond * 50)
	go func(ctx context.Context) {
		for {
			select {
			case <-ticker.C:
				d.calibrate()
			case <-ctx.Done():
				ticker.Stop()
				return
			}
		}
	}(ctx)
}

func (d *demoQueue) ProducersUp(delta uint32) error {
	if delta == 0 {
		delta = 1
	}
	if d.producersUp+delta >= d.pmax {
		return errors.New("maximum producers count reached")
	}
	c := d.producersUp
	for i := c; i < c+delta; i++ {
		go d.producers[i].produce(d.queue)
		d.producers[i].start()
		d.producersUp++
		ProducerStartMetric(d.key)
	}
	return nil
}

func (d *demoQueue) ProducersDown(delta uint32) error {
	if delta == 0 {
		delta = 1
	}
	params, _ := d.rtParams()
	if d.producersUp-delta < params.min {
		return errors.New("minimum producers count reached")
	}
	c := d.producersUp
	for i := c; i >= c-delta; i-- {
		if d.producers[i].getStatus() == statusActive {
			d.producers[i].stop()
			d.producersUp--
			ProducerStopMetric(d.key)
		}
	}
	return nil
}

func (d *demoQueue) Stop() {
	d.stop(false)
}

func (d *demoQueue) ForceStop() {
	d.stop(true)
}

func (d *demoQueue) stop(force bool) {
	c := d.producersUp
	for i := uint32(0); i < c; i++ {
		d.producers[i].stop()
		d.producersUp--
		ProducerStopMetric(d.key)
	}
	if force {
		d.queue.ForceClose()
	} else {
		d.queue.Close()
	}
	d.cancel()
}

func (d *demoQueue) calibrate() {
	var (
		params  realtimeParams
		schedID int
	)
	if params, schedID = d.rtParams(); schedID != d.schedID {
		d.schedID = schedID
		if d.pmax > params.max {
			for i := d.pmax - 1; i >= params.max; i-- {
				if d.producers[i].getStatus() == statusActive {
					if err := d.ProducersDown(1); err != nil {
						log.Println("err", err)
					}
				}
			}
		}
		if pu := d.producersUp; params.min > pu {
			target := params.min - pu
			if err := d.ProducersUp(target); err != nil {
				log.Println("err", err)
			}
		}
	}
}

func (d *demoQueue) rtParams() (params realtimeParams, schedID int) {
	if d.schedule != nil {
		var schedParams blqueue.ScheduleParams
		if schedParams, schedID = d.schedule.Get(); schedID != -1 {
			params.min, params.max = schedParams.WorkersMin, schedParams.WorkersMax
			return
		}
	}
	schedID = -1
	params.min = d.req.ProducersMin
	params.max = d.req.ProducersMax
	return
}

func (d *demoQueue) producersMaxDaily() uint32 {
	sched, conf := uint32(0), d.req.ProducersMax
	if d.schedule != nil {
		sched = d.schedule.WorkersMaxDaily()
	}
	if sched > conf {
		return sched
	}
	return conf
}

func (d *demoQueue) String() string {
	var out = &struct {
		Key             string `json:"key"`
		Queue           string `json:"queue"`
		ProducersMin    int    `json:"producers_min"`
		ProducersMax    int    `json:"producers_max"`
		ProducersIdle   int    `json:"producers_idle"`
		ProducersActive int    `json:"producers_active"`
	}{}

	out.Key = d.key
	out.Queue = "!queue"
	out.ProducersMin = int(d.req.ProducersMin)
	out.ProducersMax = int(d.pmax)
	for _, p := range d.producers {
		switch p.getStatus() {
		case statusIdle:
			out.ProducersIdle++
		case statusActive:
			out.ProducersActive++
		}
	}

	b, _ := json.Marshal(out)
	b = bytes.Replace(b, []byte(`"!queue"`), []byte(d.queue.String()), 1)

	return string(b)
}
