package request

import "github.com/koykov/demo/wirex/internal/protobuf"

type ModelContainer struct {
	User UserModelInterface
	SSP  SSPModelInterface
}

type UserModelInterface interface {
	GetUserName(id int32) string
}

type SSPModelInterface interface {
	Get(id int32) *protobuf.SSP
}
