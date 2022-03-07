package main

import (
	"encoding/json"
	"io"
	"math"
	"math/rand"
	"net/http"
	"time"

	"github.com/koykov/demo/traced"
	"github.com/koykov/demo/traced/model"
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
	defer func() { w.WriteHeader(status) }()

	time.Sleep(time.Duration(200+rand.Intn(3000)) * time.Millisecond)
	if rand.Intn(100) > 80 {
		status = http.StatusNoContent
		return
	}

	switch r.URL.Path {
	case "/v1":
		if !checkMethod(r, "POST") {
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
			Price:  float32(math.Floor(req.BF)) + rand.Float32(),
			Markup: []byte(traced.RandString(128)),
		}
		resp, _ = json.Marshal(v1)
	case "/v2":
		if !checkMethod(r, "POST") {
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
			Commission: math.Floor(req.BF) + rand.Float64(),
			Currency:   req.Cur,
			Data:       traced.RandString(64),
		}
		resp, _ = json.Marshal(v2)
	case "/v3":
		if !checkMethod(r, "GET") {
			status = http.StatusMethodNotAllowed
			return
		}
		// ...
	default:
		status = http.StatusNotFound
		return
	}

	_, _ = w.Write(resp)
}

func checkMethod(r *http.Request, must string) bool {
	return r.Method == must
}
