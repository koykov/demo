package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/koykov/batch_query"
	mw "github.com/koykov/metrics_writers/batch_query"
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

func NewBQHTTP() *BQHTTP {
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

	case r.URL.Path == "/api/v1/init":
		if bq != nil {
			resp.Status = http.StatusNotAcceptable
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println("err", err)
			resp.Status = http.StatusBadRequest
			resp.Error = err.Error()
			return
		}

		var (
			req  RequestInit
			conf batch_query.Config
		)

		err = json.Unmarshal(body, &req)
		if err != nil {
			log.Println("err", err)
			resp.Status = http.StatusBadRequest
			resp.Error = err.Error()
			return
		}
		req.MapConfig(&conf)

		conf.MetricsWriter = mw.NewPrometheusMetricsWP(key, time.Millisecond)
		conf.Logger = log.New(os.Stderr, fmt.Sprintf("queue #%s ", key), log.LstdFlags)

		var qi *batch_query.BatchQuery
		qi, _ = batch_query.New(&conf)

		q := demoBQ{
			key: key,
			bq:  qi,
			req: &req,
		}

		h.mux.Lock()
		h.pool[key] = &q
		h.mux.Unlock()

		q.Run()

		resp.Message = "success"

	case r.URL.Path == "/api/v1/producer-up" && bq != nil:
		var delta uint32
		if d := r.FormValue("delta"); len(d) > 0 {
			ud, err := strconv.ParseUint(d, 10, 32)
			if err != nil {
				log.Println("err", err)
				resp.Status = http.StatusInternalServerError
				resp.Error = err.Error()
				return
			}
			delta = uint32(ud)
		}
		if err := bq.ProducersUp(delta); err != nil {
			log.Println("err", err)
			resp.Status = http.StatusInternalServerError
			resp.Error = err.Error()
			return
		}
		resp.Message = "success"

	case r.URL.Path == "/api/v1/producer-down" && bq != nil:
		var delta uint32
		if d := r.FormValue("delta"); len(d) > 0 {
			ud, err := strconv.ParseUint(d, 10, 32)
			if err != nil {
				log.Println("err", err)
				resp.Status = http.StatusInternalServerError
				resp.Error = err.Error()
				return
			}
			delta = uint32(ud)
		}
		if err := bq.ProducersDown(delta); err != nil {
			log.Println("err", err)
			resp.Status = http.StatusInternalServerError
			resp.Error = err.Error()
			return
		}
		resp.Message = "success"

	case r.URL.Path == "/api/v1/stop":
		if bq != nil {
			bq.Stop()
		}

		h.mux.Lock()
		delete(h.pool, key)
		h.mux.Unlock()

		resp.Message = "success"

	case r.URL.Path == "/api/v1/force-stop":
		if bq != nil {
			bq.ForceStop()
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

func (h *BQHTTP) get(key string) *demoBQ {
	h.mux.RLock()
	defer h.mux.RUnlock()
	if q, ok := h.pool[key]; ok {
		return q
	}
	return nil
}
