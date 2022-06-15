package feed

import (
	"sync"

	"github.com/google/wire"
	"github.com/koykov/demo/wirex/internal/http/request"
)

var (
	hdl *Handler
	ho  sync.Once

	svc *Service
	so  sync.Once

	repo *Repository
	ro   sync.Once

	ProviderSet = wire.NewSet(
		ProvideHandler,
		ProvideService,
		ProvideRepository,

		wire.Bind(new(request.Handler), new(*Handler)),
		wire.Bind(new(request.Service), new(*Service)),
		wire.Bind(new(request.Repository), new(*Repository)),
	)
)

func ProvideHandler(svc request.Service) *Handler {
	ho.Do(func() {
		hdl = &Handler{
			svc: svc,
		}
	})

	return hdl
}

func ProvideService(repo request.Repository) *Service {
	so.Do(func() {
		svc = &Service{
			repo: repo,
		}
	})

	return svc
}

func ProvideRepository(dlv request.Deliveler, mc *request.ModelContainer) *Repository {
	ro.Do(func() {
		repo = &Repository{
			dlv:   dlv,
			model: mc,
		}
	})

	return repo
}
