package main

import (
	"time"

	"github.com/koykov/batch_query"
)

type RequestInit struct {
	QueryBatchSize       uint64        `json:"query_batch_size,omitempty"`
	QueryTimeoutInterval time.Duration `json:"query_timeout_interval,omitempty"`
	QueryCollectInterval time.Duration `json:"query_collect_interval,omitempty"`
	QueryWorkers         uint          `json:"query_workers"`
	QueryBuffer          uint64        `json:"query_buffer,omitempty"`

	ProducersMin  uint32 `json:"producers_min"`
	ProducersMax  uint32 `json:"producers_max"`
	ProducerDelay uint32 `json:"producer_delay,omitempty"`

	Aerospike *struct {
		Host            string        `json:"host"`
		Port            int           `json:"port"`
		Namespace       string        `json:"namespace"`
		SetName         string        `json:"set_name"`
		Bins            []string      `json:"bins"`
		ReadTimeoutNS   time.Duration `json:"read_timeout_ns"`
		TotalTimeoutNS  time.Duration `json:"total_timeout_ns"`
		SocketTimeoutNS time.Duration `json:"socket_timeout_ns"`
		MaxRetries      int           `json:"max_retries"`
	} `json:"aerospike"`
}

func (r *RequestInit) MapConfig(conf *batch_query.Config) {
	conf.BatchSize = r.QueryBatchSize
	conf.CollectInterval = r.QueryCollectInterval
	conf.TimeoutInterval = r.QueryTimeoutInterval
	conf.Workers = r.QueryWorkers
	conf.Buffer = r.QueryBuffer
}
