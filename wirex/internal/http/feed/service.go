package feed

import (
	"github.com/koykov/demo/wirex/internal/http/request"
)

type Service struct {
	repo request.Repository
}

func (s Service) GetRequest() (*request.Inner, error) {
	return nil, nil
}
