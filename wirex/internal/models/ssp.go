package models

import "github.com/koykov/demo/wirex/internal/protobuf"

type SSPModel struct{}

func (m SSPModel) Get(id int32) *protobuf.Source {
	return &protobuf.Source{
		ID:   id,
		Name: "foobar",
	}
}
