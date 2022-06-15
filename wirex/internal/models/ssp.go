package models

import "github.com/koykov/demo/wirex/internal/protobuf"

type SSPModel struct{}

func (m SSPModel) Get(id int32) *protobuf.SSP {
	return &protobuf.SSP{
		ID:   id,
		Name: "foobar",
	}
}
