package main

import (
	"encoding/json"
	"os"

	"github.com/koykov/traceID"
)

type Config struct {
	Listen struct {
		AppPort uint `json:"appPort"`
		CbPort  uint `json:"cbPort"`
		PbPort  uint `json:"pbPort"`
	} `json:"listen"`
	Clients []uint `json:"clients"`

	Broadcaster traceID.BroadcasterConfig `json:"broadcaster"`
}

func (c *Config) LoadFrom(file string) error {
	b, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, c)
}
