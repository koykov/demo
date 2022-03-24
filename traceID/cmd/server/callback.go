package main

import (
	"net/http"
	"regexp"

	td "github.com/koykov/demo/traceID"
	"github.com/koykov/demo/traceID/model"
	"github.com/koykov/fastconv"
	"github.com/koykov/traceID"
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

	ttx.SetID(req.TraceID).SetService("cbw")
	ttx.Info("income /cb request").
		Var("method", r.Method).
		Var("url", r.URL)
	ttx.Debug("request body").
		Var("encoded", m[1]).
		Var("decoded", req)

	// ...
}
