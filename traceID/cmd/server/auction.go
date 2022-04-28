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
	resp model.Response
	err  error
}

type streamRE chan re

func Auction(ttx *traceID.Ctx, req *model.Request) (resp *model.Response, err error) {
	var pool []CV
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
					winner = &re.resp
				}
			}
		}
	}(ctx, ttx, stream)

	var wg sync.WaitGroup
	for i := 0; i < len(pool); i++ {
		wg.Add(1)
		go execReq(ttx, &pool[i], req, stream, &wg)
	}
	wg.Wait()

	close(stream)

	if winner == nil {
		ttx.Warn("no winner")
		return
	}

	ttx.Info("auction winner").
		Var("winner", winner)

	if req.BF > 0 && winner.Bid < req.BF {
		err = ErrBidfloorFail
		ttx.Warn("bidfloor check failed").
			Var("bidfloor", req.BF).
			Var("bid", winner.Bid)
		return
	}
	if req.BC > 0 && winner.Bid > req.BC {
		err = ErrBidceilFail
		ttx.Warn("bidceil check failed").
			Var("bidceil", req.BC).
			Var("bid", winner.Bid)
		return
	}

	resp = winner

	return
}

func execReq(ttx *traceID.Ctx, cv *CV, req *model.Request, stream streamRE, wg *sync.WaitGroup) {
	var (
		resp re
		hr   *http.Response
		buf  []byte
	)

	tth := ttx.AcquireThread()
	defer func() {
		ttx.ReleaseThread(tth)
		wg.Done()
	}()

	defer func() { stream <- resp }()
	switch cv.Version {
	case "v1":
		b := req.ToV1()
		tth.Info("send request").
			Var("addr", cv.Client).
			Var("version", cv.Version).
			Var("body", fastconv.B2S(b))
		if hr, resp.err = http.Post(cv.Client+"/"+cv.Version, "application/json", bytes.NewBuffer(b)); resp.err != nil {
			tth.Error("request failed").Err(resp.err)
			return
		}
		tth.Debug("request v1 done").
			Var("code", hr.StatusCode).
			Var("len", hr.ContentLength)
		buf, resp.err = io.ReadAll(hr.Body)
		if resp.err != nil {
			tth.Error("body read failed").Err(resp.err)
			return
		}
		tth.Debug("response v1 body").
			Var("body", string(buf))
		if resp.err = resp.resp.FromV1(buf); resp.err != nil {
			tth.Error("body decoding failed").
				Var("version", cv.Version).
				Err(resp.err)
			return
		}
		tth.Debug("decoded response v1").
			Var("decoded", resp.resp)
	case "v2":
		b := req.ToV2()
		tth.Info("send request").
			Var("addr", cv.Client).
			Var("version", cv.Version).
			Var("body", fastconv.B2S(b))
		if hr, resp.err = http.Post(cv.Client+"/"+cv.Version, "application/json", bytes.NewBuffer(b)); resp.err != nil {
			tth.Error("request failed").Err(resp.err)
			return
		}
		tth.Debug("request v2 done").
			Var("code", hr.StatusCode).
			Var("len", hr.ContentLength)
		if buf, resp.err = io.ReadAll(hr.Body); resp.err != nil {
			tth.Error("body read failed").Err(resp.err)
			return
		}
		tth.Debug("response v2 body").
			Var("body", string(buf))
		if resp.err = resp.resp.FromV2(buf); resp.err != nil {
			tth.Error("body decoding failed").
				Var("version", cv.Version).
				Err(resp.err)
			return
		}
		tth.Debug("decoded response v2").
			Var("decoded", resp.resp)
	case "v3":
		b := req.ToV3()
		tth.Info("send request").
			Var("addr", cv.Client).
			Var("version", cv.Version).
			Var("url", fastconv.B2S(b))
		if hr, resp.err = http.Get(cv.Client + string(b)); resp.err != nil {
			tth.Error("request failed").Err(resp.err)
			return
		}
		tth.Debug("request v3 done").
			Var("code", hr.StatusCode).
			Var("len", hr.ContentLength)
		if buf, resp.err = io.ReadAll(hr.Body); resp.err != nil {
			tth.Error("body read failed").Err(resp.err)
			return
		}
		tth.Debug("response v3 body").
			Var("body", string(buf))
		if resp.err = resp.resp.FromV3(buf); resp.err != nil {
			tth.Error("body decoding failed").
				Var("version", cv.Version).
				Err(resp.err)
			return
		}
		tth.Debug("decoded response v3").
			Var("decoded", resp.resp)
	}
	return
}
