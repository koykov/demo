package main

import (
	"errors"
	"math/rand"
	"time"
)

type Worker struct {
	delay    uint32
	deadline bool
	failRate float32
}

func NewWorker(delay uint32, allowDeadline bool, failRate float32) *Worker {
	d := &Worker{
		delay:    delay,
		deadline: allowDeadline,
		failRate: failRate,
	}
	if d.delay == 0 {
		d.delay = 75
	}
	return d
}

func (d *Worker) Do(_ interface{}) error {
	var delta int
	if d.deadline {
		delta = rand.Intn(int(d.delay)) - int(d.delay/2)
	}
	delay := time.Duration(d.delay) + time.Duration(delta)
	time.Sleep(delay)
	if d.failRate > 0 {
		r := rand.Float32()
		if r < d.failRate {
			return errArtificialFail
		}
	}
	return nil
}

var errArtificialFail = errors.New("artificial fail")
