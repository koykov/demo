package models

import "github.com/koykov/demo/wirex/internal/protobuf"

type SourceModel struct{}

func (m SourceModel) Get(id int32) *protobuf.Source {
	return &protobuf.Source{
		ID:   id,
		Name: "foobar",
	}
}
