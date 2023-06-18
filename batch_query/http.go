package main

import (
	"encoding/json"
	"net/http"
	"sync"
)

type BQHTTP struct {
	mux  sync.RWMutex
	pool map[string]*demoBQ

	allow400 map[string]bool
	allow404 map[string]bool
}

type BQResponse struct {
	Status  int    `json:"status,omitempty"`
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}

func NewBQHTTP(hport, pport int) *BQHTTP {
	h := &BQHTTP{
		pool: make(map[string]*demoBQ),
		allow400: map[string]bool{
			"/api/v1/ping": true,
			"/api/v1/list": true,
		},
		allow404: map[string]bool{
			"/api/v1/init": true,
			"/api/v1/ping": true,
			"/api/v1/list": true,
		},
	}
	return h
}

func (h *BQHTTP) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		key  string
		bq   *demoBQ
		resp BQResponse
	)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, X-Auth-Token")

	defer func() {
		w.WriteHeader(resp.Status)
		b, _ := json.Marshal(resp)
		_, _ = w.Write(b)
	}()

	resp.Status = http.StatusOK

	if key = r.FormValue("key"); len(key) == 0 && !h.allow400[r.URL.Path] {
		resp.Status = http.StatusBadRequest
		return
	}
	if bq = h.get(key); bq == nil && !h.allow404[r.URL.Path] {
		resp.Status = http.StatusNotFound
		return
	}

	switch {
	case r.URL.Path == "/api/v1/ping":
		resp.Message = "pong"
	default:
		resp.Status = http.StatusNotFound
		return
	}
}

func (h *BQHTTP) get(key string) *demoBQ {
	h.mux.RLock()
	defer h.mux.RUnlock()
	if q, ok := h.pool[key]; ok {
		return q
	}
	return nil
}
