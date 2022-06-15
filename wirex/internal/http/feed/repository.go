package feed

import (
	"github.com/koykov/demo/wirex/internal/http/request"
)

type Repository struct {
	dlv   request.Deliveler
	model *request.ModelContainer
	proc  request.ProcFnContainer
}

func (r Repository) GetRequest() (*request.Origin, error) {
	req := request.Origin{
		Value: r.dlv.RequestURI(),
		Body:  r.dlv.PostBody(),
	}
	return &req, nil
}
