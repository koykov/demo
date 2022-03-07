package model

import "encoding/json"

type Response struct {
	Bid  float64 `json:"bid"`
	Cur  string  `json:"cur"`
	Mkup string  `json:"mkup"`
}

type ResponseV1 struct {
	Price  float32 `json:"price"`
	Markup []byte  `json:"markup"`
}

type ResponseV2 struct {
	Commission float64 `json:"commission"`
	Currency   string  `json:"currency"`
	Data       string  `json:"data"`
}

type ResponseV3 struct {
	A float32 `json:"a"`
	B string  `json:"b"`
	C string  `json:"c"`
}

func (r *Response) FromV1(p []byte) (err error) {
	var v1 ResponseV1
	if err = json.Unmarshal(p, &v1); err != nil {
		return
	}
	r.Bid = float64(v1.Price)
	r.Mkup = string(v1.Markup)
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
	return
}
