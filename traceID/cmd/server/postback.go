package main

import (
	"math/rand"
	"net/http"
	"regexp"
	"time"

	td "github.com/koykov/demo/traceID"
	"github.com/koykov/demo/traceID/model"
	"github.com/koykov/fastconv"
	"github.com/koykov/traceID"
	"github.com/koykov/traceID/marshaller"
)

type PostbackHTTP struct{}

var (
	rePB = regexp.MustCompile(`/pb/(.*)`)
)

func (h *PostbackHTTP) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, X-Auth-Token")

	var out []byte
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
	if m = rePB.FindStringSubmatch(r.URL.Path); len(m) == 0 {
		status = http.StatusBadRequest
		return
	}

	var req model.PBRequest
	if err := req.Unmarshal([]byte(m[1])); err != nil {
		status = http.StatusBadRequest
		return
	}
	if len(req.TraceID) == 0 {
		status = http.StatusBadRequest
		return
	}

	ttx.SetID(req.TraceID).SetServiceWithStage("pbworker", req.UniqID)
	ttx.Info("income /pb request").
		Var("method", r.Method).
		Var("url", r.URL)
	ttx.Debug("request body").
		Var("encoded", m[1]).
		Var("decoded", req)

	time.Sleep(100 + time.Duration(rand.Intn(300)))

	body := td.RandString(20)
	ttx.Info("PB response").
		Var("body", body)

	out = []byte(body)
}
