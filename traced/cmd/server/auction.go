package main

import (
	"bytes"
	"io"
	"net/http"
	"sync"

	"github.com/koykov/demo/traced/model"
)

type re struct {
	resp *model.Response
	err  error
}

type streamRE chan re

func Auction(req *model.Request) (resp *model.Response, err error) {
	var pool []cv
	if pool, err = filterClients(req); err != nil {
		return
	}
	stream := make(streamRE, len(pool))
	var wg sync.WaitGroup
	for i := 0; i < len(pool); i++ {
		wg.Add(1)
		go execReq(&pool[i], req, stream)
	}
	wg.Wait()
	return
}

func execReq(cv *cv, req *model.Request, stream streamRE) {
	var (
		resp re
		hr   *http.Response
		buf  []byte
	)
	defer func() { stream <- resp }()
	switch cv.v {
	case "v1":
		b := req.ToV1()
		if hr, resp.err = http.Post(cv.c, "application/json", bytes.NewBuffer(b)); resp.err != nil {
			return
		}
		buf, resp.err = io.ReadAll(hr.Body)
		if resp.err != nil {
			return
		}
		if resp.err = resp.resp.FromV1(buf); resp.err != nil {
			return
		}
	case "v2":
		b := req.ToV1()
		if hr, resp.err = http.Post(cv.c, "application/json", bytes.NewBuffer(b)); resp.err != nil {
			return
		}
		if buf, resp.err = io.ReadAll(hr.Body); resp.err != nil {
			return
		}
		if resp.err = resp.resp.FromV2(buf); resp.err != nil {
			return
		}
	case "v3":
		b := req.ToV1()
		if hr, resp.err = http.Get(cv.c + string(b)); resp.err != nil {
			return
		}
		if buf, resp.err = io.ReadAll(hr.Body); resp.err != nil {
			return
		}
		if resp.err = resp.resp.FromV3(buf); resp.err != nil {
			return
		}
	}
	return
}
