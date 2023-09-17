package main

import (
	"bufio"
	"math/rand"
	"os"
	"strconv"
)

const maxKey = 10_000

type keysRepo struct {
	buf []int64
}

func (r *keysRepo) load(keysPath string) error {
	f, err := os.Open(keysPath)
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
	if len(r.buf) == 0 {
		return rand.Int63n(maxKey)
	}
	return r.buf[rand.Intn(len(r.buf))]
}
