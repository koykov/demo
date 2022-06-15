package request

import "github.com/koykov/demo/wirex/internal/protobuf"

type ModelContainer struct {
	User   UserModelInterface
	Source SourceModelInterface
}

type UserModelInterface interface {
	GetUserName(id int32) string
}

type SourceModelInterface interface {
	Get(id int32) *protobuf.Source
}
