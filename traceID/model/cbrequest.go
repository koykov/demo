package model

import (
	"bytes"
	"encoding/base32"
	"encoding/json"
)

type CBRequest struct {
	Bid     float64 `json:"bid"`
	Cur     string  `json:"cur"`
	PB      string  `json:"pb"`
	TraceID string  `json:"trace_id"`
}

func (r CBRequest) Marshal() ([]byte, error) {
	jb, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	buf := make([]byte, base32.StdEncoding.EncodedLen(len(jb)))
	base32.StdEncoding.Encode(buf, jb)
	return bytes.ToLower(buf), nil
}

func (r *CBRequest) Unmarshal(p []byte) error {
	p = bytes.ToUpper(p)
	buf := make([]byte, base32.StdEncoding.DecodedLen(len(p)))
	n, err := base32.StdEncoding.Decode(buf, p)
	if err != nil {
		return err
	}
	return json.Unmarshal(buf[:n], r)
}
