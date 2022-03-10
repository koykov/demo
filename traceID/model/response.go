package model

import "encoding/json"

type Response struct {
	Bid     float64 `json:"bid"`
	Cur     string  `json:"cur"`
	Mkup    string  `json:"mkup"`
	TraceID string  `json:"trace_id"`
}

type ResponseV1 struct {
	Price   float32 `json:"price"`
	Markup  []byte  `json:"markup"`
	TraceID string  `json:"trace_id"`
}

type ResponseV2 struct {
	Commission float64 `json:"commission"`
	Currency   string  `json:"currency"`
	Data       string  `json:"data"`
	TraceID    string  `json:"trace_id"`
}

type ResponseV3 struct {
	A       float32 `json:"a"`
	B       string  `json:"b"`
	C       string  `json:"c"`
	TraceID string  `json:"trace_id"`
}

func (r *Response) FromV1(p []byte) (err error) {
	var v1 ResponseV1
	if err = json.Unmarshal(p, &v1); err != nil {
		return
	}
	r.Bid = float64(v1.Price)
	r.Cur = "USD"
	r.Mkup = string(v1.Markup)
	r.TraceID = v1.TraceID
	return
}

func (r *Response) FromV2(p []byte) (err error) {
	var v2 ResponseV2
	if err = json.Unmarshal(p, &v2); err != nil {
		return
	}
	r.Bid = v2.Commission
	r.Mkup = v2.Data
	r.Cur = v2.Currency
	r.TraceID = v2.TraceID
	return
}

func (r *Response) FromV3(p []byte) (err error) {
	var v3 ResponseV3
	if err = json.Unmarshal(p, &v3); err != nil {
		return
	}
	r.Bid = float64(v3.A)
	r.Mkup = v3.B
	r.Cur = v3.C
	r.TraceID = v3.TraceID
	return
}
