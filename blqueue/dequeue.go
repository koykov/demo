package main

import "time"

type Worker struct {
	delay uint32
}

func NewWorker(delay uint32) *Worker {
	d := &Worker{delay: delay}
	if d.delay == 0 {
		d.delay = 75
	}
	return d
}

func (d *Worker) Do(_ interface{}) error {
	time.Sleep(time.Duration(d.delay) * time.Nanosecond)
	return nil
}
