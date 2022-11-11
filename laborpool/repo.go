package main

type RequestInit struct {
	Size          uint    `json:"size"`
	PensionFactor float32 `json:"pension_factor"`
	ProducersMin  uint32  `json:"producers_min"`
	ProducersMax  uint32  `json:"producers_max"`
	ProducerDelay uint32  `json:"producer_delay,omitempty"`
	WorkerDelay   uint32  `json:"worker_delay,omitempty"`
}
