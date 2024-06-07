package main

import (
	"encoding/binary"
	"math"

	"github.com/koykov/byteconv"
)

type entryEncoder[T any] struct{}

func (entryEncoder[T]) Encode(dst []byte, t T) ([]byte, int, error) {
	off := len(dst)
	e := any(t).(entry)
	dst = binary.LittleEndian.AppendUint64(dst, uint64(e.c))
	dst = binary.LittleEndian.AppendUint64(dst, uint64(e.i))
	dst = binary.LittleEndian.AppendUint64(dst, e.u)
	f := math.Float64bits(e.f)
	dst = binary.LittleEndian.AppendUint64(dst, f)
	dst = binary.LittleEndian.AppendUint32(dst, uint32(len(e.s)))
	dst = append(dst, e.s...)
	dst = binary.LittleEndian.AppendUint32(dst, uint32(len(e.b)))
	dst = append(dst, e.b...)
	return dst, len(dst) - off, nil
}

type entryDecoder[T any] struct{}

func (entryDecoder[T]) Decode(x *T, p []byte) error {
	e := any(*x).(entry)
	e.c = int(binary.LittleEndian.Uint64(p))
	p = p[8:]
	e.i = int64(binary.LittleEndian.Uint64(p))
	p = p[8:]
	e.u = binary.LittleEndian.Uint64(p)
	p = p[8:]
	fu := binary.LittleEndian.Uint64(p)
	e.f = math.Float64frombits(fu)
	p = p[8:]
	l := binary.LittleEndian.Uint32(p)
	p = p[4:]
	e.s = byteconv.B2S(p[:l])
	p = p[l:]
	l = binary.LittleEndian.Uint32(p)
	p = p[4:]
	e.b = p[:l]
	return nil
}
