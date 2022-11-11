package main

import (
	"encoding/json"
	"errors"

	"github.com/koykov/laborpool"
)

type demoPool struct {
	key  string
	pool *laborpool.Pool
	req  *RequestInit

	producersUp uint32
	producers   []*producer
}

func (d *demoPool) Run() {
	d.producers = make([]*producer, d.req.ProducersMax)
	for i := 0; i < int(d.req.ProducersMax); i++ {
		d.producers[i] = makeProducer(uint32(i), d.req.ProducerDelay)
	}
	for i := 0; i < int(d.req.ProducersMin); i++ {
		go d.producers[i].produce(d.pool)
		d.producers[i].start()
	}
	d.producersUp = d.req.ProducersMin
	ProducersInitMetric(d.key, d.producersUp, d.req.ProducersMax-d.producersUp)
}

func (d *demoPool) ProducersUp(delta uint32) error {
	if delta == 0 {
		delta = 1
	}
	if d.producersUp+delta >= d.req.ProducersMax {
		return errors.New("maximum producers count reached")
	}
	c := d.producersUp
	for i := c; i < c+delta; i++ {
		go d.producers[i].produce(d.pool)
		d.producers[i].start()
		d.producersUp++
		ProducerStartMetric(d.key)
	}
	return nil
}

func (d *demoPool) ProducersDown(delta uint32) error {
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

func (d *demoPool) Stop() {
	c := d.producersUp
	for i := uint32(0); i < c; i++ {
		d.producers[i].stop()
		d.producersUp--
		ProducerStopMetric(d.key)
	}
}

func (d *demoPool) String() string {
	var out = &struct {
		Key             string  `json:"key"`
		Size            uint    `json:"size"`
		PensionFactor   float32 `json:"pension_factor"`
		ProducersMin    int     `json:"producers_min"`
		ProducersMax    int     `json:"producers_max"`
		ProducersIdle   int     `json:"producers_idle"`
		ProducersActive int     `json:"producers_active"`
	}{}

	out.Key = d.key
	out.Size = d.req.Size
	out.PensionFactor = d.req.PensionFactor
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

	return string(b)
}
