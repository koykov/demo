package main

import (
	"errors"
	"math/rand"

	"github.com/koykov/demo/traceID/model"
)

type CV struct {
	Client  string
	Version string
}

var (
	cpool []string
	vp    = []string{"v1", "v2", "v3"}
)

func filterClients(req *model.Request) (r []CV, err error) {
	for i := 0; i < len(cpool); i++ {
		if rand.Intn(100) > 80 {
			continue
		}
		v := vp[rand.Intn(3)]
		r = append(r, CV{
			Client:  cpool[i],
			Version: v,
		})
	}
	if len(r) == 0 {
		err = errors.New("no clients available")
	}
	return
}
