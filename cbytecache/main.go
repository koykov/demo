package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/koykov/cbytecache"
	"github.com/koykov/hash/fnv"
	metrics "github.com/koykov/metric_writers/cbytecache"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	dataSize = 1e6
	writers  = 10
	readers  = 10
	testDur  = 30 * time.Second
	sleepDur = 100 * time.Microsecond
)

var (
	config = cbytecache.Config{
		HashFn:        fnv.Hash64aString,
		Buckets:       4,
		Expire:        5 * time.Minute,
		Vacuum:        300 * time.Minute,
		MaxSize:       1 * cbytecache.Gigabyte,
		MetricsWriter: metrics.NewPrometheusMetrics("demo1GB"),
		Logger:        log.New(os.Stdout, "", log.LstdFlags),
	}

	exmpl = [][]byte{
		[]byte(`{"firstName":"John","lastName":"Smith","isAlive":true,"age":27,"address":{"streetAddress":"21 2nd Street","city":"New York","state":"NY","postalCode":"10021-3100"},"phoneNumbers":[{"type":"home","number":"212 555-1234"},{"type":"office","number":"646 555-4567"},{"type":"mobile","number":"123 456-7890"}],"children":[],"spouse":null}`),
		[]byte(`{"$schema":"http://json-schema.org/schema#","title":"Product","type":"object","required":["id","name","price"],"properties":{"id":{"type":"number","description":"Product identifier"},"name":{"type":"string","description":"Name of the product"},"price":{"type":"number","minimum":0},"tags":{"type":"array","items":{"type":"string"}},"stock":{"type":"object","properties":{"warehouse":{"type":"number"},"retail":{"type":"number"}}}}}`),
		[]byte(`{"id":1,"name":"Foo","price":123,"tags":["Bar","Eek"],"stock":{"warehouse":300,"retail":20}}`),
		[]byte(`{"first name":"John","last name":"Smith","age":25,"address":{"street address":"21 2nd Street","city":"New York","state":"NY","postal code":"10021"},"phone numbers":[{"type":"home","number":"212 555-1234"},{"type":"fax","number":"646 555-4567"}],"sex":{"type":"male"}}`),
		[]byte(`{"fruit":"Apple","size":"Large","color":"Red"}`),
		[]byte(`{"quiz":{"sport":{"q1":{"question":"Which one is correct team name in NBA?","options":["New York Bulls","Los Angeles Kings","Golden State Warriros","Huston Rocket"],"answer":"Huston Rocket"}},"maths":{"q1":{"question":"5 + 7 = ?","options":["10","11","12","13"],"answer":"12"},"q2":{"question":"12 - 8 = ?","options":["1","2","3","4"],"answer":"4"}}}}`),
	}
	data = make(map[string][]byte, dataSize)

	pport = flag.Int("pport", 8081, "Prometheus port")
)

func init() {
	flag.Parse()

	eln := int32(len(exmpl))
	for i := 0; i < dataSize; i++ {
		key := fmt.Sprintf("key%d", i)
		payload := exmpl[rand.Int31n(eln)]
		data[key] = payload
	}
}

func main() {
	paddr := fmt.Sprintf(":%d", *pport)
	go func() {
		// registered metrics endpoint
		http.Handle("/metrics", promhttp.Handler())
		log.Printf("Start Prometheus server on address %s/metrics\n", paddr)
		if err := http.ListenAndServe(paddr, nil); err != nil {
			log.Fatal(err)
		}
	}()

	c, err := cbytecache.NewCByteCache(&config)
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup

	for i := 0; i < writers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			d := time.After(testDur)
			for {
				select {
				case <-d:
					return
				default:
					time.Sleep(sleepDur)
					key := fmt.Sprintf("key%d", rand.Int31n(dataSize))
					payload := data[key]
					if err := c.Set(key, payload); err != nil {
						log.Println("err", err)
					}
				}
			}
		}()
	}

	for i := 0; i < readers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			d := time.After(testDur)
			for {
				select {
				case <-d:
					return
				default:
					time.Sleep(sleepDur)
					var (
						dst []byte
						err error
					)
					key := fmt.Sprintf("key%d", rand.Int31n(dataSize))
					payload := data[key]
					if dst, err = c.Get(key); err != nil && err != cbytecache.ErrNotFound {
						log.Println("err", err)
						continue
					}
					if len(dst) > 0 && !bytes.Equal(dst, payload) {
						log.Println("mismatch: key", key, "got", string(dst), "need", string(payload))
					}
				}
			}
		}()
	}

	wg.Wait()

	log.Println("done")
}
