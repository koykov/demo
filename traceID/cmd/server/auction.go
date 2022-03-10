package main

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"sync"

	"github.com/koykov/demo/traceID/model"
	"github.com/koykov/fastconv"
	"github.com/koykov/traceID"
)

type re struct {
	resp *model.Response
	err  error
}

type streamRE chan re

func Auction(ttx *traceID.Ctx, req *model.Request) (resp *model.Response, err error) {
	var pool []cv
	if pool, err = filterClients(req); err != nil {
		ttx.Error("build clients pool failed").Err(err)
		return
	}
	ttx.Info("clients pool").Var("list", pool)

	stream := make(streamRE, len(pool))
	var (
		winner *model.Response
		maxBid float64
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func(ctx context.Context, ttx *traceID.Ctx, stream streamRE) {
		var c int
		tth := ttx.AcquireThread()
		defer func() {
			ttx.Debug("reviewed {count} responses").
				Var("count", c).
				Var("max bid", maxBid)
			ttx.ReleaseThread(tth)
		}()
		for {
			select {
			case <-ctx.Done():
				return
			case re, ok := <-stream:
				if !ok {
					return
				}
				if re.err != nil {
					continue
				}
				c++
				if re.resp.Bid > maxBid {
					winner = re.resp
				}
			}
		}
	}(ctx, ttx, stream)

	var wg sync.WaitGroup
	for i := 0; i < len(pool); i++ {
		wg.Add(1)
		go execReq(ttx, &pool[i], req, stream)
	}
	wg.Wait()

	close(stream)

	// todo consider bid floor/ceil

	ttx.Info("auction winner").
		Var("winner", winner)

	return
}

func execReq(ttx *traceID.Ctx, cv *cv, req *model.Request, stream streamRE) {
	var (
		resp re
		hr   *http.Response
		buf  []byte
	)

	tth := ttx.AcquireThread()
	defer ttx.ReleaseThread(tth)

	defer func() { stream <- resp }()
	switch cv.v {
	case "v1":
		b := req.ToV1()
		tth.Info("send request").
			Var("addr", cv.c).
			Var("version", cv.v).
			Var("body", fastconv.B2S(b))
		if hr, resp.err = http.Post(cv.c, "application/json", bytes.NewBuffer(b)); resp.err != nil {
			tth.Error("request failed").
				Var("code", hr.StatusCode).
				Err(resp.err)
			return
		}
		buf, resp.err = io.ReadAll(hr.Body)
		if resp.err != nil {
			tth.Error("body read failed").Err(resp.err)
			return
		}
		if resp.err = resp.resp.FromV1(buf); resp.err != nil {
			tth.Error("body decoding failed").
				Var("version", cv.v).
				Err(resp.err)
			return
		}
	case "v2":
		b := req.ToV2()
		tth.Info("send request").
			Var("addr", cv.c).
			Var("version", cv.v).
			Var("body", fastconv.B2S(b))
		if hr, resp.err = http.Post(cv.c, "application/json", bytes.NewBuffer(b)); resp.err != nil {
			tth.Error("request failed").
				Var("code", hr.StatusCode).
				Err(resp.err)
			return
		}
		if buf, resp.err = io.ReadAll(hr.Body); resp.err != nil {
			tth.Error("body read failed").Err(resp.err)
			return
		}
		if resp.err = resp.resp.FromV2(buf); resp.err != nil {
			tth.Error("body decoding failed").
				Var("version", cv.v).
				Err(resp.err)
			return
		}
	case "v3":
		b := req.ToV3()
		tth.Info("send request").
			Var("addr", cv.c).
			Var("version", cv.v).
			Var("url", fastconv.B2S(b))
		if hr, resp.err = http.Get(cv.c + string(b)); resp.err != nil {
			tth.Error("request failed").
				Var("code", hr.StatusCode).
				Err(resp.err)
			return
		}
		if buf, resp.err = io.ReadAll(hr.Body); resp.err != nil {
			tth.Error("body read failed").Err(resp.err)
			return
		}
		if resp.err = resp.resp.FromV3(buf); resp.err != nil {
			tth.Error("body decoding failed").
				Var("version", cv.v).
				Err(resp.err)
			return
		}
	}
	return
}
