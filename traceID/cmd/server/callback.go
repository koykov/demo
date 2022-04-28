package main

import (
	"context"
	"io"
	"net/http"
	"regexp"
	"time"

	td "github.com/koykov/demo/traceID"
	"github.com/koykov/demo/traceID/model"
	"github.com/koykov/fastconv"
	"github.com/koykov/traceID"
	"github.com/koykov/traceID/marshaller"
)

type CallbackHTTP struct{}

var (
	reCB = regexp.MustCompile(`/cb/(.*)`)
)

func (h *CallbackHTTP) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, X-Auth-Token")

	var (
		out []byte
	)
	status := http.StatusOK
	ttx := traceID.AcquireCtx()
	ttx.SetLogger(logger).SetMarshaller(marshaller.JSON{})
	defer func() {
		w.WriteHeader(status)
		_, _ = w.Write(out)

		ttx.Info("response").
			Var("status", status).
			Var("body", fastconv.B2S(out))

		_ = ttx.Flush()
		traceID.ReleaseCtx(ttx)
	}()

	if !td.CheckMethod(r, "GET") {
		status = http.StatusMethodNotAllowed
		return
	}

	var m []string
	if m = reCB.FindStringSubmatch(r.URL.Path); len(m) == 0 {
		status = http.StatusBadRequest
		return
	}

	var req model.CBRequest
	if err := req.Unmarshal([]byte(m[1])); err != nil {
		status = http.StatusBadRequest
		return
	}
	if len(req.TraceID) == 0 {
		status = http.StatusBadRequest
		return
	}

	ttx.SetID(req.TraceID).SetService("cbworker")
	ttx.Info("income /cb request").
		Var("method", r.Method).
		Var("url", r.URL)
	ttx.Debug("request body").
		Var("encoded", m[1]).
		Var("decoded", req)

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*300)
	defer cancel()
	hreq, err := http.NewRequestWithContext(ctx, "GET", req.PB, nil)
	if err != nil {
		ttx.Error("PB request fail").
			Var("stage", "build").
			Err(err)
		status = http.StatusInternalServerError
		return
	}
	hres, err := http.DefaultClient.Do(hreq)
	if err != nil {
		ttx.Error("PB request fail").
			Var("stage", "exec").
			Err(err)
		status = http.StatusInternalServerError
		return
	}
	body, err := io.ReadAll(hres.Body)
	defer func() { _ = hres.Body.Close() }()
	if err != nil {
		ttx.Error("PB request fail").
			Var("stage", "read").
			Err(err)
		status = http.StatusInternalServerError
		return
	}
	ttx.Info("PB response").
		Var("code", hres.StatusCode).
		Var("body", string(body))

	out = body
}
