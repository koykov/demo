package feed

import (
	"github.com/koykov/demo/wirex/internal/http/request"
	"github.com/valyala/fasthttp"
)

type Handler struct {
	svc request.Service
}

func (h Handler) Process() fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		// ...
	}
}
