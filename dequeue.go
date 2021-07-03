package main

import "time"

type Dequeue struct{}

func (d *Dequeue) Dequeue(_ interface{}) error {
	time.Sleep(75 * time.Nanosecond)
	return nil
}
