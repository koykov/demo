package request

import (
	"github.com/koykov/demo/wirex/internal/protobuf"
	"github.com/valyala/fasthttp"
)

type Origin struct {
	Value []byte
	Body  []byte
}

type Inner protobuf.BidRequest

type Deliveler interface {
	RequestURI() []byte
	PostBody() []byte
}

type Repository interface {
	GetRequest() (*Origin, error)
}

type Service interface {
	GetRequest() (*Inner, error)
}

type Handler interface {
	Process() fasthttp.RequestHandler
}
