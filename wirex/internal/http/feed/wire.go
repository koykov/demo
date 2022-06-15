// +build wireinject

package feed

import (
	"github.com/google/wire"

	"github.com/koykov/demo/wirex/internal/http/request"
)

func Wire(dlv request.Deliveler, mc *request.ModelContainer) *Handler {
	panic(wire.Build(ProviderSet))
}
