package main

import (
	"encoding/json"
	"net/http"
	"sync"
)

type CacheHTTP struct {
	mux  sync.RWMutex
	pool map[string]*demoCache

	hport, pport int
}

type CacheResponse struct {
	Status  int    `json:"status,omitempty"`
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}

func NewCacheHTTP(hport, pport int) *CacheHTTP {
	h := &CacheHTTP{
		pool:  make(map[string]*demoCache),
		hport: hport,
		pport: pport,
	}
	return h
}

func (h *CacheHTTP) get(key string) *demoCache {
	h.mux.RLock()
	defer h.mux.RUnlock()
	if c, ok := h.pool[key]; ok {
		return c
	}
	return nil
}

func (h *CacheHTTP) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		// key  string
		// c    *demoCache
		resp CacheResponse
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

	switch {
	case r.URL.Path == "/api/v1/ping":
		resp.Message = "pong"
	default:
		resp.Status = http.StatusNotFound
		return
	}
}
