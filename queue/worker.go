package main

import (
	"math/rand"
	"time"
)

type Worker struct {
	delay    uint32
	deadline bool
}

func NewWorker(delay uint32, allowDeadline bool) *Worker {
	d := &Worker{delay: delay, deadline: allowDeadline}
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
	return nil
}
