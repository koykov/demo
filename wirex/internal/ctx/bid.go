package ctx

import "github.com/koykov/demo/wirex/internal/http/request"

type BidCtx struct {
	Models struct {
		User request.UserModelInterface
		SSP  request.SSPModelInterface
	}
}
