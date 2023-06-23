package main

import (
	"time"

	"github.com/koykov/batch_query"
)

type RequestInit struct {
	QueryChunkSize       uint64        `json:"query_chunk_size,omitempty"`
	QueryCollectInterval time.Duration `json:"query_collect_interval,omitempty"`
	QueryWorkers         uint          `json:"query_workers"`
	QueryBuffer          uint64        `json:"query_buffer,omitempty"`

	ProducersMin  uint32 `json:"producers_min"`
	ProducersMax  uint32 `json:"producers_max"`
	ProducerDelay uint32 `json:"producer_delay,omitempty"`
}

func (r *RequestInit) MapConfig(conf *batch_query.Config) {
	conf.ChunkSize = r.QueryChunkSize
	conf.CollectInterval = r.QueryCollectInterval
	conf.Workers = r.QueryWorkers
	conf.Buffer = r.QueryBuffer
}
