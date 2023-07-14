package main

import (
	"bufio"
	"math/rand"
	"os"
	"strconv"
)

type keysRepo struct {
	buf []int64
}

func (r *keysRepo) load() error {
	f, err := os.Open("batch_query/aerospike.txt")
	if err != nil {
		return err
	}
	defer f.Close()
	scnr := bufio.NewScanner(f)
	for scnr.Scan() {
		raw := scnr.Text()
		key, err := strconv.Atoi(raw)
		if err != nil {
			return err
		}
		r.buf = append(r.buf, int64(key))
	}
	return nil
}

func (r *keysRepo) get() int64 {
	return r.buf[rand.Intn(len(r.buf))]
}
