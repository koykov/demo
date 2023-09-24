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
		KeysPath        string        `json:"keys_path"`
		Host            string        `json:"host"`
		Port            int           `json:"port"`
		Instances       uint          `json:"instances"`
		Namespace       string        `json:"namespace"`
		SetName         string        `json:"set_name"`
		Bins            []string      `json:"bins"`
		ReadTimeoutNS   time.Duration `json:"read_timeout_ns"`
		TotalTimeoutNS  time.Duration `json:"total_timeout_ns"`
		SocketTimeoutNS time.Duration `json:"socket_timeout_ns"`
		MaxRetries      int           `json:"max_retries"`
	} `json:"aerospike"`
	Mysql *DBConfig `json:"mysql"`
	Pgsql *DBConfig `json:"pgsql"`
	Redis *struct {
		KeysPath           string        `json:"keys_path"`
		DSN                string        `json:"dsn"`
		Network            string        `json:"network"`
		Addr               string        `json:"addr"`
		Pass               string        `json:"pass"`
		DB                 int           `json:"db"`
		MaxRetries         int           `json:"max_retries"`
		MinRetryBackoff    time.Duration `json:"min_retry_backoff"`
		MaxRetryBackoff    time.Duration `json:"max_retry_backoff"`
		DialTimeout        time.Duration `json:"dial_timeout"`
		ReadTimeout        time.Duration `json:"read_timeout"`
		WriteTimeout       time.Duration `json:"write_timeout"`
		PoolSize           int           `json:"pool_size"`
		MinIdleConns       int           `json:"min_idle_conns"`
		MaxConnAge         time.Duration `json:"max_conn_age"`
		PoolTimeout        time.Duration `json:"pool_timeout"`
		IdleTimeout        time.Duration `json:"idle_timeout"`
		IdleCheckFrequency time.Duration `json:"idle_check_frequency"`
	} `json:"redis"`
}

type DBConfig struct {
	DSN      string `json:"dsn"`
	Addr     string `json:"addr"`
	Host     string `json:"host"`
	Port     uint   `json:"port"`
	User     string `json:"user"`
	Pass     string `json:"pass"`
	Protocol string `json:"protocol"`
	DBName   string `json:"db_name"`
	DDL      string `json:"ddl"`
	DML      bool   `json:"dml"`
}

func (r *RequestInit) MapConfig(conf *batch_query.Config) {
	conf.BatchSize = r.QueryBatchSize
	conf.CollectInterval = r.QueryCollectInterval
	conf.TimeoutInterval = r.QueryTimeoutInterval
	conf.Workers = r.QueryWorkers
	conf.Buffer = r.QueryBuffer
}
