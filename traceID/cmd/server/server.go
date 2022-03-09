package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/koykov/demo/traceID"
	"github.com/koykov/demo/traceID/model"
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
	defer func() {
		w.WriteHeader(status)
		_, _ = w.Write(out)
	}()

	switch r.URL.Path {
	case "/v1":
		if !traceID.CheckMethod(r, "POST") {
			status = http.StatusMethodNotAllowed
			return
		}
		body, err = io.ReadAll(r.Body)
		if err != nil {
			status = http.StatusInternalServerError
			return
		}
		if err = req.FromV1(body); err != nil {
			status = http.StatusBadRequest
			return
		}
		if resp, err = Auction(&req); err != nil {
			status = http.StatusInternalServerError
			return
		}

		out, _ = json.Marshal(resp)

	case "/v2":
		if !traceID.CheckMethod(r, "POST") {
			status = http.StatusMethodNotAllowed
			return
		}
		body, err = io.ReadAll(r.Body)
		if err != nil {
			status = http.StatusInternalServerError
			return
		}
		if err = req.FromV2(body); err != nil {
			status = http.StatusBadRequest
			return
		}
		if resp, err = Auction(&req); err != nil {
			status = http.StatusInternalServerError
			return
		}

		out, _ = json.Marshal(resp)

	case "/v3":
		if !traceID.CheckMethod(r, "GET") {
			status = http.StatusMethodNotAllowed
			return
		}
		if err = req.FromV3([]byte(r.URL.RawQuery)); err != nil {
			status = http.StatusBadRequest
			return
		}
		if resp, err = Auction(&req); err != nil {
			status = http.StatusInternalServerError
			return
		}

		out, _ = json.Marshal(resp)
	}
}
