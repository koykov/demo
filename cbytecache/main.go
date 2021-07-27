package main

import (
	"log"
	"os"
	"time"

	"github.com/koykov/cbytecache"
	"github.com/koykov/fastconv"
)

var (
	config = cbytecache.Config{
		HashFn:        fastconv.Fnv64aString,
		Shards:        4,
		Expire:        30 * time.Second,
		Vacuum:        300 * time.Second,
		MaxSize:       512 * cbytecache.Kilobyte,
		MetricsWriter: nil,
		Logger:        log.New(os.Stdout, "", log.LstdFlags),
	}
)

func main() {
	c, err := cbytecache.NewCByteCache(config)
	if err != nil {
		log.Fatal(err)
	}
	err = c.Set("foo", []byte("bar"))
	log.Println(err)
}
