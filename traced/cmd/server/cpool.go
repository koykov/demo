package main

import (
	"errors"
	"math/rand"

	"github.com/koykov/demo/traced/model"
)

type cv struct {
	c string
	v string
}

var (
	cpool []string
	vp    = []string{"v1", "v2", "v3"}
)

func filterClients(req *model.Request) (r []cv, err error) {
	for i := 0; i < len(cpool); i++ {
		if rand.Intn(100) > 80 {
			continue
		}
		v := vp[rand.Intn(3)]
		r = append(r, cv{
			c: cpool[i],
			v: v,
		})
	}
	if len(r) == 0 {
		err = errors.New("no clients available")
	}
	return
}
