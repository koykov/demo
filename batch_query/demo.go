package main

import (
	"context"
	"io"
	"time"
)

type demoBQ struct {
	key string
	req *RequestInit

	producersUp uint32
	producers   []*producer

	cancel context.CancelFunc
}

func (d *demoBQ) Run() {
	d.producers = make([]*producer, d.req.ProducersMax)
	for i := 0; i < int(d.req.ProducersMax); i++ {
		d.producers[i] = makeProducer(uint32(i), d.req.ProducerDelay, d.req.AllowDeadline, d.req.WorkerDelay)
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

func (d *demoBQ) ProducersUp(delta uint32) error {
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

func (d *demoBQ) ProducersDown(delta uint32) error {
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

func (d *demoBQ) Stop() {
	d.stop(false)
}

func (d *demoBQ) ForceStop() {
	d.stop(true)
}

func (d *demoBQ) stop(force bool) {
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
