package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/koykov/cbytecache"
	"github.com/koykov/clock"
	"github.com/koykov/hash/fnv"
	metrics "github.com/koykov/metrics_writers/cbytecache"
)

type CacheHTTP struct {
	mux  sync.RWMutex
	pool map[string]*demoCache

	hport, pport int

	allow400 map[string]bool
	allow404 map[string]bool
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
		key  string
		c    *demoCache
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

	if key = r.FormValue("key"); len(key) == 0 && !h.allow400[r.URL.Path] {
		resp.Status = http.StatusBadRequest
		return
	}
	if c = h.get(key); c == nil && !h.allow404[r.URL.Path] {
		resp.Status = http.StatusNotFound
		return
	}

	switch {
	case r.URL.Path == "/api/v1/ping":
		resp.Message = "pong"

	case r.URL.Path == "/api/v1/list":
		buf := bytes.Buffer{}
		buf.WriteByte('[')
		c := 0
		for _, a := range h.pool {
			if c > 0 {
				buf.WriteByte(',')
			}
			_, _ = buf.WriteString(a.String())
			c++
		}
		buf.WriteByte(']')
		resp.Message = buf.String()

	case r.URL.Path == "/api/v1/status" && c != nil:
		resp.Message = c.String()

	case r.URL.Path == "/api/v1/init":
		if c != nil {
			resp.Status = http.StatusNotAcceptable
			return
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("err", err)
			resp.Status = http.StatusBadRequest
			resp.Error = err.Error()
			return
		}

		var (
			req  RequestInit
			conf cbytecache.Config
		)

		err = json.Unmarshal(body, &req)
		if err != nil {
			log.Println("err", err)
			resp.Status = http.StatusBadRequest
			resp.Error = err.Error()
			return
		}

		req.MapConfig(&conf)

		conf.Hasher = fnv.Hasher{}
		conf.MetricsWriter = metrics.NewPrometheusMetrics(key)
		conf.Logger = log.New(os.Stderr, "", log.LstdFlags)
		conf.Clock = clock.NewClock()

		ci, err := cbytecache.New(&conf)
		if err != nil {
			log.Println("err", err)
			resp.Status = http.StatusInternalServerError
			resp.Error = err.Error()
			return
		}

		c := demoCache{
			key:     key,
			config:  &conf,
			rawReq:  &req,
			cache:   ci,
			writers: req.Writers,
			readers: req.Readers,
		}

		h.mux.Lock()
		h.pool[key] = &c
		h.mux.Unlock()

		c.Run()

		resp.Message = "success"

	case r.URL.Path == "/api/v1/stop":
		if c != nil {
			c.Stop()
		}

		h.mux.Lock()
		delete(h.pool, key)
		h.mux.Unlock()

		resp.Message = "success"

	default:
		resp.Status = http.StatusNotFound
		return
	}
}
