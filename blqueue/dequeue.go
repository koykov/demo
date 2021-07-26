package main

import "time"

type Dequeue struct {
	delay uint32
}

func NewDequeue(delay uint32) *Dequeue {
	d := &Dequeue{delay: delay}
	if d.delay == 0 {
		d.delay = 75
	}
	return d
}

func (d *Dequeue) Dequeue(_ interface{}) error {
	time.Sleep(time.Duration(d.delay) * time.Nanosecond)
	return nil
}
