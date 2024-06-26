package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	as "github.com/aerospike/aerospike-client-go/v7"
	"github.com/go-sql-driver/mysql"
	"github.com/koykov/batch_query"
	"github.com/koykov/batch_query/mods/aerospike"
	bqsql "github.com/koykov/batch_query/mods/sql"
	"github.com/koykov/demo/batch_query/ddl"
	mw "github.com/koykov/metrics_writers/batch_query"
	_ "github.com/lib/pq"
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

		switch {
		case req.Aerospike != nil:
			if err = h.initAS(&req, &conf); err != nil {
				log.Println("err", err)
				resp.Status = http.StatusInternalServerError
				resp.Error = err.Error()
			}
		case req.Mysql != nil || req.Pgsql != nil:
			if err = h.initSQL(&req, &conf); err != nil {
				log.Println("err", err)
				resp.Status = http.StatusInternalServerError
				resp.Error = err.Error()
			}
		default:
			log.Println(fmt.Errorf("no mod config provided"))
			resp.Status = http.StatusBadRequest
			resp.Error = err.Error()
			return
		}
		conf.MetricsWriter = mw.NewPrometheusMetricsWP(key, time.Millisecond)
		conf.Logger = log.New(os.Stderr, fmt.Sprintf("query #%s ", key), log.LstdFlags)

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

func (h *BQHTTP) initAS(req *RequestInit, conf *batch_query.Config) error {
	asc := req.Aerospike

	if err := krepo.load(asc.KeysPath); err != nil {
		return err
	}

	readPolicy := as.NewClientPolicy()
	readPolicy.Timeout = asc.ReadTimeoutNS

	batchPolicy := as.NewBatchPolicy()
	batchPolicy.TotalTimeout = asc.TotalTimeoutNS
	batchPolicy.SocketTimeout = asc.SocketTimeoutNS
	batchPolicy.MaxRetries = asc.MaxRetries

	inst := asc.Instances
	if inst < 2 {
		client, err := as.NewClientWithPolicy(readPolicy, asc.Host, asc.Port)
		if err != nil {
			return err
		}
		conf.Batcher = aerospike.Batcher{
			Namespace: asc.Namespace,
			SetName:   asc.SetName,
			Bins:      asc.Bins,
			Policy:    batchPolicy,
			Client:    client,
		}
	} else {
		clients := make([]*as.Client, 0, inst)
		for i := uint(0); i < inst; i++ {
			client, err := as.NewClientWithPolicy(readPolicy, asc.Host, asc.Port)
			if err != nil {
				return err
			}
			clients = append(clients, client)
		}
		conf.Batcher = aerospike.MCBatcher{
			Namespace: asc.Namespace,
			SetName:   asc.SetName,
			Bins:      asc.Bins,
			Policy:    batchPolicy,
			Clients:   clients,
		}
	}

	return nil
}

func (h *BQHTTP) initSQL(req *RequestInit, conf *batch_query.Config) error {
	var (
		dsn, typ, pfx string
		dbc           *DBConfig
		pt            bqsql.PlaceholderType
	)
	switch {
	case req.Mysql != nil:
		if dsn = req.Mysql.DSN; len(dsn) == 0 {
			cfg := mysql.Config{
				User:   req.Mysql.User,
				Passwd: req.Mysql.Pass,
				Net:    req.Mysql.Protocol,
				Addr:   req.Mysql.Addr,
				DBName: req.Mysql.DBName,
			}
			dsn = cfg.FormatDSN()
		}
		typ = "mysql"
		dbc = req.Mysql
		pt = bqsql.PlaceholderMySQL
	case req.Pgsql != nil:
		if dsn = req.Pgsql.DSN; len(dsn) == 0 {
			dsn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
				req.Pgsql.Host, req.Pgsql.Port, req.Pgsql.User, req.Pgsql.Pass, req.Pgsql.DBName)
		}
		typ = "postgres"
		dbc = req.Pgsql
		pt = bqsql.PlaceholderPgSQL
		pfx = "bq."
	}

	db, err := sql.Open(typ, dsn)
	if err != nil {
		return err
	}

	if len(dbc.DDL) > 0 {
		if err = ddl.ApplyDDL(db, dbc.DDL); err != nil {
			return err
		}
	}
	if dbc.DML {
		if err = ddl.ApplyDML(db, maxKey, pt); err != nil {
			return err
		}
	}

	rec := &SQLRecord{}
	conf.Batcher = bqsql.Batcher{
		DB:             db,
		Query:          fmt.Sprintf("select id, name, status, bio, balance from %susers where id in (::args::)", pfx),
		QueryFormatter: bqsql.MacrosQueryFormatter{PlaceholderType: pt},
		RecordScanner:  rec,
		RecordMatcher:  rec,
	}

	return nil
}

func (h *BQHTTP) get(key string) *demoBQ {
	h.mux.RLock()
	defer h.mux.RUnlock()
	if q, ok := h.pool[key]; ok {
		return q
	}
	return nil
}
