package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/koykov/blqueue"
	metrics "github.com/koykov/metrics_writers/blqueue"
)

type QueueHTTP struct {
	mux  sync.RWMutex
	pool map[string]*demoQueue

	hport, pport int

	allow400 map[string]bool
	allow404 map[string]bool
}

type QueueResponse struct {
	Status  int    `json:"status,omitempty"`
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}

func NewQueueHTTP(hport, pport int) *QueueHTTP {
	h := &QueueHTTP{
		pool:  make(map[string]*demoQueue),
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

func (h *QueueHTTP) get(key string) *demoQueue {
	h.mux.RLock()
	defer h.mux.RUnlock()
	if q, ok := h.pool[key]; ok {
		return q
	}
	return nil
}

func (h *QueueHTTP) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		key  string
		q    *demoQueue
		resp QueueResponse
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
	if q = h.get(key); q == nil && !h.allow404[r.URL.Path] {
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

	case r.URL.Path == "/api/v1/status" && q != nil:
		resp.Message = q.String()

	case r.URL.Path == "/api/v1/init":
		if q != nil {
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
			conf blqueue.Config
		)

		err = json.Unmarshal(body, &req)
		if err != nil {
			log.Println("err", err)
			resp.Status = http.StatusBadRequest
			resp.Error = err.Error()
			return
		}
		conf.Key = key
		req.MapConfig(&conf)

		// sched := blqueue.NewSchedule()
		// now := time.Now()
		// now1 := now.Add(time.Minute)
		// now2 := now.Add(time.Minute*2)
		// now3 := now.Add(time.Minute*3)
		// now4 := now.Add(time.Minute*4)
		// now5 := now.Add(time.Minute*5)
		// now6 := now.Add(time.Minute*6)
		// now7 := now.Add(time.Minute*7)
		// now8 := now.Add(time.Minute*8)
		// t1 := fmt.Sprintf("%02d:%02d:%02d", now1.Hour(), now1.Minute(), now1.Second())
		// t2 := fmt.Sprintf("%02d:%02d:%02d", now2.Hour(), now2.Minute(), now2.Second())
		// t3 := fmt.Sprintf("%02d:%02d:%02d", now3.Hour(), now3.Minute(), now3.Second())
		// t4 := fmt.Sprintf("%02d:%02d:%02d", now4.Hour(), now4.Minute(), now4.Second())
		// t5 := fmt.Sprintf("%02d:%02d:%02d", now5.Hour(), now5.Minute(), now5.Second())
		// t6 := fmt.Sprintf("%02d:%02d:%02d", now6.Hour(), now6.Minute(), now6.Second())
		// t7 := fmt.Sprintf("%02d:%02d:%02d", now7.Hour(), now7.Minute(), now7.Second())
		// t8 := fmt.Sprintf("%02d:%02d:%02d", now8.Hour(), now8.Minute(), now8.Second())
		// _ = sched.AddRange(t1 + "-" + t2, blqueue.ScheduleParams{WorkersMin: 4, WorkersMax: 16})
		// _ = sched.AddRange(t3 + "-" + t4, blqueue.ScheduleParams{WorkersMin: 1, WorkersMax: 32})
		// _ = sched.AddRange(t5 + "-" + t6, blqueue.ScheduleParams{WorkersMin: 10, WorkersMax: 16})
		// _ = sched.AddRange(t7 + "-" + t8, blqueue.ScheduleParams{WorkersMin: 1, WorkersMax: 20})
		// log.Println(sched.String())
		// conf.Schedule = sched

		conf.MetricsWriter = metrics.NewPrometheusMetrics()
		conf.Dequeuer = NewDequeue(req.WorkerDelay)
		if req.AllowLeak {
			conf.DLQ = &blqueue.DummyDLQ{}
		}

		conf.Logger = log.New(os.Stderr, "", log.LstdFlags)

		qi, _ := blqueue.New(&conf)

		q := demoQueue{
			key:           key,
			queue:         qi,
			producersMin:  req.ProducersMin,
			producersMax:  req.ProducersMax,
			producerDelay: req.ProducerDelay,
		}

		h.mux.Lock()
		h.pool[key] = &q
		h.mux.Unlock()

		q.Run()

		resp.Message = "success"

	case r.URL.Path == "/api/v1/producer-up" && q != nil:
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
		if err := q.ProducerUp(delta); err != nil {
			log.Println("err", err)
			resp.Status = http.StatusInternalServerError
			resp.Error = err.Error()
			return
		}
		resp.Message = "success"

	case r.URL.Path == "/api/v1/producer-down" && q != nil:
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
		if err := q.ProducerDown(delta); err != nil {
			log.Println("err", err)
			resp.Status = http.StatusInternalServerError
			resp.Error = err.Error()
			return
		}
		resp.Message = "success"

	case r.URL.Path == "/api/v1/stop":
		if q != nil {
			q.Stop()
		}

		h.mux.Lock()
		delete(h.pool, key)
		h.mux.Unlock()

		resp.Message = "success"

	case r.URL.Path == "/api/v1/force-stop":
		if q != nil {
			q.ForceStop()
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
