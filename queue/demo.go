package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"time"

	"github.com/koykov/dlqdump"
	"github.com/koykov/queue"
)

type demoQueue struct {
	key   string
	queue *queue.Queue
	req   *RequestInit
	dlq   queue.Enqueuer
	rst   *dlqdump.Restorer

	producersUp uint32
	producers   []*producer

	schedule *queue.Schedule
	schedID  int

	cancel context.CancelFunc
}

func (d *demoQueue) Run() {
	d.producers = make([]*producer, d.req.ProducersMax)
	for i := 0; i < int(d.req.ProducersMax); i++ {
		d.producers[i] = makeProducer(uint32(i), d.req.ProducerDelay)
	}
	for i := 0; i < int(d.req.ProducersMin); i++ {
		go d.producers[i].produce(d.queue)
		d.producers[i].start()
	}
	d.producersUp = d.req.ProducersMin
	ProducersInitMetric(d.key, d.producersUp, d.req.ProducersMax-d.producersUp)

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
	if d.producersUp+delta-1 >= d.req.ProducersMax {
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
	if d.producersUp-delta < d.req.ProducersMin {
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
		_ = d.queue.ForceClose()
		if d.rst != nil {
			_ = d.rst.ForceClose()
		}
	} else {
		_ = d.queue.Close()
		if d.rst != nil {
			_ = d.rst.Close()
		}
	}
	if d.dlq != nil {
		inst := any(d.dlq).(io.Closer)
		_ = inst.Close()
	}
	d.cancel()
}

func (d *demoQueue) calibrate() {
	var (
		rtSize  uint32
		schedID int
	)
	if rtSize, schedID = d.rtSize(); schedID != d.schedID {
		d.schedID = schedID
		pu := d.producersUp
		if rtSize > pu {
			target := rtSize - pu
			if err := d.ProducersUp(target); err != nil {
				log.Println("err", err)
			}
		} else if rtSize < pu {
			target := pu - rtSize
			if err := d.ProducersDown(target); err != nil {
				log.Println("err", err)
			}
		}
	}
}

func (d *demoQueue) rtSize() (size uint32, schedID int) {
	if d.schedule != nil {
		var schedParams queue.ScheduleParams
		if schedParams, schedID = d.schedule.Get(); schedID != -1 {
			size = schedParams.WorkersMin
			return
		}
	}
	schedID = -1
	size = d.req.ProducersMin
	return
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
	out.ProducersMax = int(d.req.ProducersMax)
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
