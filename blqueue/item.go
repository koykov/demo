package main

import (
	"encoding/binary"
	"io"
)

type Item struct {
	Header  uint32
	Payload uint64
}

func (i Item) Size() int {
	return 12
}

func (i Item) MarshalTo(p []byte) (int, error) {
	if len(p) < i.Size() {
		return 0, io.ErrShortBuffer
	}
	binary.LittleEndian.PutUint32(p[:4], i.Header)
	binary.LittleEndian.PutUint64(p[4:], i.Payload)
	return i.Size(), nil
}
