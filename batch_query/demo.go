package main

import (
	"errors"

	"github.com/koykov/batch_query"
)

type demoBQ struct {
	key string
	bq  *batch_query.BatchQuery
	req *RequestInit

	producersUp uint32
	producers   []*producer
}

func (d *demoBQ) Run() {
	d.producers = make([]*producer, d.req.ProducersMax)
	for i := 0; i < int(d.req.ProducersMax); i++ {
		d.producers[i] = makeProducer(uint32(i))
	}
	for i := 0; i < int(d.req.ProducersMin); i++ {
		go d.producers[i].produce(d.bq)
		d.producers[i].start()
	}
	d.producersUp = d.req.ProducersMin
	ProducersInitMetric(d.key, d.producersUp, d.req.ProducersMax-d.producersUp)
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
		go d.producers[i].produce(d.bq)
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
	if c >= d.req.ProducersMax {
		c = d.req.ProducersMax - 1
	}
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
		_ = d.bq.ForceClose()
	} else {
		_ = d.bq.Close()
	}
}
