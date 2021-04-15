package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"

	"github.com/koykov/blqueue"
)

type demoQueue struct {
	key   string
	queue *blqueue.Queue

	hport, pport int

	producersMin,
	producersMax,
	producersUp uint32
	producers []*producer

	stats [500]Stat
}

type Stat struct {
	Qsize, Qleak,
	Wactive, Wsleep, Widle,
	Pactive, Pidle int
}

var (
	reQsize   = regexp.MustCompile(`queue_size{queue=".*"} (\d+)`)
	reQleak   = regexp.MustCompile(`queue_leak{queue=".*"} (\d+)`)
	reWactive = regexp.MustCompile(`queue_workers_active{queue=".*"} (\d+)`)
	reWsleep  = regexp.MustCompile(`queue_workers_sleep{queue=".*"} (\d+)`)
	reWidle   = regexp.MustCompile(`queue_workers_idle{queue=".*"} (\d+)`)
	rePactive = regexp.MustCompile(`queue_producers_active{queue=".*"} (\d+)`)
	rePidle   = regexp.MustCompile(`queue_producers_idle{queue=".*"} (\d+)`)
)

func (d *demoQueue) Run() {
	d.producers = make([]*producer, d.producersMax)
	for i := 0; i < int(d.producersMax); i++ {
		d.producers[i] = makeProducer(uint32(i))
	}
	for i := 0; i < int(d.producersMin); i++ {
		go d.producers[i].produce(d.queue)
		d.producers[i].start()
	}
	d.producersUp = d.producersMin

	producerActive.WithLabelValues(d.key).Add(float64(d.producersUp))
	producerIdle.WithLabelValues(d.key).Add(float64(d.producersMax - d.producersUp))
}

func (d *demoQueue) ProducerUp(delta uint32) error {
	if delta == 0 {
		delta = 1
	}
	if d.producersUp+delta >= d.producersMax {
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

func (d *demoQueue) ProducerDown(delta uint32) error {
	if delta == 0 {
		delta = 1
	}
	if d.producersUp-delta < d.producersMin {
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
	c := d.producersUp
	for i := uint32(0); i < c; i++ {
		d.producers[i].stop()
		d.producersUp--
		ProducerStopMetric(d.key)
	}
	d.queue.Close()
}

func (d *demoQueue) String() string {
	var out = &struct {
		Key             string    `json:"key"`
		Queue           string    `json:"queue"`
		ProducersMin    int       `json:"producers_min"`
		ProducersMax    int       `json:"producers_max"`
		ProducersIdle   int       `json:"producers_idle"`
		ProducersActive int       `json:"producers_active"`
		Stats           [500]Stat `json:"stats"`
	}{}

	out.Key = d.key
	out.Queue = "!queue"
	out.ProducersMin = int(d.producersMin)
	out.ProducersMax = int(d.producersMax)
	for _, p := range d.producers {
		switch p.getStatus() {
		case statusIdle:
			out.ProducersIdle++
		case statusActive:
			out.ProducersActive++
		}
	}

	resp, _ := http.Get("http://localhost:" + strconv.Itoa(d.pport) + "/metrics")
	defer func() { _ = resp.Body.Close() }()
	contents, _ := ioutil.ReadAll(resp.Body)

	var (
		qsize, qleak,
		wactive, wsleep, widle,
		pactive, pidle int64
	)
	if m := reQsize.FindSubmatch(contents); m != nil {
		qsize, _ = strconv.ParseInt(string(m[1]), 10, 64)
	}
	if m := reQleak.FindSubmatch(contents); m != nil {
		qleak, _ = strconv.ParseInt(string(m[1]), 10, 64)
	}
	if m := reWactive.FindSubmatch(contents); m != nil {
		wactive, _ = strconv.ParseInt(string(m[1]), 10, 64)
	}
	if m := reWsleep.FindSubmatch(contents); m != nil {
		wsleep, _ = strconv.ParseInt(string(m[1]), 10, 64)
	}
	if m := reWidle.FindSubmatch(contents); m != nil {
		widle, _ = strconv.ParseInt(string(m[1]), 10, 64)
	}
	if m := rePactive.FindSubmatch(contents); m != nil {
		pactive, _ = strconv.ParseInt(string(m[1]), 10, 64)
	}
	if m := rePidle.FindSubmatch(contents); m != nil {
		pidle, _ = strconv.ParseInt(string(m[1]), 10, 64)
	}

	copy(d.stats[0:], d.stats[1:500])
	d.stats[499] = Stat{
		Qsize:   int(qsize),
		Qleak:   int(qleak),
		Wactive: int(wactive),
		Wsleep:  int(wsleep),
		Widle:   int(widle),
		Pactive: int(pactive),
		Pidle:   int(pidle),
	}
	out.Stats = d.stats

	b, _ := json.Marshal(out)
	b = bytes.Replace(b, []byte(`"!queue"`), []byte(d.queue.String()), 1)

	return string(b)
}
