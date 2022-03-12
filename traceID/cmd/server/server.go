package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	td "github.com/koykov/demo/traceID"
	"github.com/koykov/demo/traceID/model"
	"github.com/koykov/fastconv"
	"github.com/koykov/traceID"
)

type ServerHTTP struct{}

var (
	logger = log.New(os.Stdout, "", log.LstdFlags)
)

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
	ttx.SetLogger(logger)
	defer func() {
		w.WriteHeader(status)
		_, _ = w.Write(out)

		ttx.Info("response").
			Var("status", status).
			Var("body", fastconv.B2S(out))

		_ = ttx.Flush()
		traceID.ReleaseCtx(ttx)
	}()

	var id string
	if v := r.URL.Query()["traceID"]; len(v) > 0 {
		id = v[0]
	}
	if len(id) == 0 {
		status = http.StatusBadRequest
		return
	}

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
	}
	req.TraceID = id

	if resp, err = Auction(ttx, &req); err != nil {
		ttx.Error("auction failed").Err(err)
		status = http.StatusInternalServerError
		return
	}
	out, _ = json.Marshal(resp)
}
