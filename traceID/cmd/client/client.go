package main

import (
	"encoding/json"
	"io"
	"math"
	"math/rand"
	"net/http"
	"time"

	"github.com/koykov/demo/traceID"
	"github.com/koykov/demo/traceID/model"
)

type ClientHTTP struct{}

func (h *ClientHTTP) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, X-Auth-Token")

	var (
		req  model.Request
		resp []byte
	)
	status := http.StatusOK
	defer func() {
		w.WriteHeader(status)
		_, _ = w.Write(resp)
	}()

	time.Sleep(time.Duration(200+rand.Intn(3000)) * time.Millisecond)
	if rand.Intn(100) > 80 {
		status = http.StatusNoContent
		return
	}

	switch r.URL.Path {
	case "/v1":
		if !traceID.CheckMethod(r, "POST") {
			status = http.StatusMethodNotAllowed
			return
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			status = http.StatusInternalServerError
			return
		}
		if err = req.FromV1(body); err != nil {
			status = http.StatusBadRequest
			return
		}
		v1 := model.ResponseV1{
			Price:   float32(randomizeBid(req.BF, req.BC)),
			Markup:  []byte(traceID.RandString(128)),
			TraceID: req.TraceID,
		}
		resp, _ = json.Marshal(v1)
	case "/v2":
		if !traceID.CheckMethod(r, "POST") {
			status = http.StatusMethodNotAllowed
			return
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			status = http.StatusInternalServerError
			return
		}
		if err = req.FromV2(body); err != nil {
			status = http.StatusBadRequest
			return
		}
		v2 := model.ResponseV2{
			Commission: randomizeBid(req.BF, req.BC),
			Currency:   req.Cur,
			Data:       traceID.RandString(64),
			TraceID:    req.TraceID,
		}
		resp, _ = json.Marshal(v2)
	case "/v3":
		if !traceID.CheckMethod(r, "GET") {
			status = http.StatusMethodNotAllowed
			return
		}
		if err := req.FromV3([]byte(r.URL.RawQuery)); err != nil {
			status = http.StatusBadRequest
			return
		}
		v3 := model.ResponseV3{
			A:       float32(randomizeBid(req.BF, req.BC)),
			B:       traceID.RandString(32),
			C:       req.Cur,
			TraceID: req.TraceID,
		}
		resp, _ = json.Marshal(v3)
	default:
		status = http.StatusNotFound
		return
	}
}

func randomizeBid(bf, bc float64) float64 {
	b := math.Floor(bf) + rand.Float64()
	if rand.Intn(100) > 95 {
		b = bf - 0.001
	}
	if rand.Intn(100) < 5 {
		b = bc + 0.001
	}
	return b
}
