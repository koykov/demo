package model

import (
	"bytes"
	"encoding/base32"
	"encoding/binary"
	"errors"
	"math"
)

type PBRequest struct {
	Commission float32 `json:"commission"`
	Cur        string  `json:"cur"`
	TraceID    string  `json:"trace_id"`
	UniqID     string  `json:"uniq_id,omitempty"`
}

func (r PBRequest) Marshal() ([]byte, error) {
	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.LittleEndian, r.Commission); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.LittleEndian, uint16(len(r.Cur))); err != nil {
		return nil, err
	}
	if _, err := buf.WriteString(r.Cur); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.LittleEndian, uint16(len(r.TraceID))); err != nil {
		return nil, err
	}
	if _, err := buf.WriteString(r.TraceID); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.LittleEndian, uint16(len(r.UniqID))); err != nil {
		return nil, err
	}
	if _, err := buf.WriteString(r.UniqID); err != nil {
		return nil, err
	}
	buf1 := make([]byte, base32.StdEncoding.EncodedLen(buf.Len()))
	base32.StdEncoding.Encode(buf1, buf.Bytes())
	return bytes.ToLower(buf1), nil
}

func (r *PBRequest) Unmarshal(p []byte) error {
	p = bytes.ToUpper(p)
	buf := make([]byte, base32.StdEncoding.DecodedLen(len(p)))
	n, err := base32.StdEncoding.Decode(buf, p)
	if err != nil {
		return err
	}
	p = buf[:n]

	if len(p) < 4 {
		return ErrPacketTooShort
	}
	r.Commission = math.Float32frombits(binary.LittleEndian.Uint32(p[:4]))
	p = p[4:]

	var cl, tl uint16
	if len(p) < 2 {
		return ErrPacketTooShort
	}
	cl = binary.LittleEndian.Uint16(p[:2])
	p = p[2:]
	if uint16(len(p)) < cl {
		return ErrPacketTooShort
	}
	r.Cur = string(p[:cl])
	p = p[cl:]
	if len(p) < 2 {
		return ErrPacketTooShort
	}
	tl = binary.LittleEndian.Uint16(p[:2])
	p = p[2:]
	if uint16(len(p)) < tl {
		return ErrPacketTooShort
	}
	r.TraceID = string(p[:tl])
	p = p[tl:]
	if len(p) < 2 {
		return ErrPacketTooShort
	}
	tl = binary.LittleEndian.Uint16(p[:2])
	p = p[2:]
	if uint16(len(p)) < tl {
		return ErrPacketTooShort
	}
	r.UniqID = string(p[:tl])
	return nil
}

var ErrPacketTooShort = errors.New("packet too short")
