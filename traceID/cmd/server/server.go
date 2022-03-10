package main

import (
	"encoding/json"
	"io"
	"net/http"

	td "github.com/koykov/demo/traceID"
	"github.com/koykov/demo/traceID/model"
	"github.com/koykov/fastconv"
	"github.com/koykov/traceID"
)

type ServerHTTP struct{}

func (h *ServerHTTP) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, X-Auth-Token")

	var (
		req  model.Request
		resp *model.Response
		out  []byte
		body []byte
		err  error
	)
	status := http.StatusOK
	ttx := traceID.AcquireCtx()
	defer func() {
		w.WriteHeader(status)
		_, _ = w.Write(out)

		ttx.Info("response").
			Var("status", status).
			Var("body", fastconv.B2S(out))

		_ = ttx.Flush()
		traceID.ReleaseCtx(ttx)
	}()

	switch r.URL.Path {
	case "/v1":
		ttx.Info("income /v1 request").
			Var("method", r.Method).
			Var("url", r.URL)
		if !td.CheckMethod(r, "POST") {
			ttx.Error("method mismatch").
				Var("need", "POST").
				Var("got", r.Method)
			status = http.StatusMethodNotAllowed
			return
		}
		body, err = io.ReadAll(r.Body)
		if err != nil {
			ttx.Error("body read error").Err(err)
			status = http.StatusInternalServerError
			return
		}
		ttx.Info("request body").Var("body", fastconv.B2S(body))
		if err = req.FromV1(body); err != nil {
			ttx.Error("body decoding from v1 failed").Err(err)
			status = http.StatusBadRequest
			return
		}
		if resp, err = Auction(ttx, &req); err != nil {
			ttx.Error("auction failed").Err(err)
			status = http.StatusInternalServerError
			return
		}

		out, _ = json.Marshal(resp)

	case "/v2":
		ttx.Info("income /v1 request").
			Var("method", r.Method).
			Var("url", r.URL)
		if !td.CheckMethod(r, "POST") {
			ttx.Error("method mismatch").
				Var("need", "POST").
				Var("got", r.Method)
			status = http.StatusMethodNotAllowed
			return
		}
		body, err = io.ReadAll(r.Body)
		if err != nil {
			ttx.Error("body read error").Err(err)
			status = http.StatusInternalServerError
			return
		}
		if err = req.FromV2(body); err != nil {
			ttx.Error("body decoding from v2 failed").Err(err)
			status = http.StatusBadRequest
			return
		}
		if resp, err = Auction(ttx, &req); err != nil {
			ttx.Error("auction failed").Err(err)
			status = http.StatusInternalServerError
			return
		}

		out, _ = json.Marshal(resp)

	case "/v3":
		ttx.Info("income /v1 request").
			Var("method", r.Method).
			Var("url", r.URL)
		if !td.CheckMethod(r, "GET") {
			ttx.Error("method mismatch").
				Var("need", "GET").
				Var("got", r.Method)
			status = http.StatusMethodNotAllowed
			return
		}
		if err = req.FromV3([]byte(r.URL.RawQuery)); err != nil {
			ttx.Error("body decoding from v3 failed").Err(err)
			status = http.StatusBadRequest
			return
		}
		if resp, err = Auction(ttx, &req); err != nil {
			ttx.Error("auction failed").Err(err)
			status = http.StatusInternalServerError
			return
		}

		out, _ = json.Marshal(resp)
	}
}
