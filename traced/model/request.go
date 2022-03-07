package model

import (
	"encoding/json"
	"net/url"
	"strconv"
)

type Request struct {
	BF    float64 `json:"bf"`
	BC    float64 `json:"bc"`
	Limit uint    `json:"limit"`
	UID   string  `json:"uid"`
	Cur   string  `json:"cur"`
}

type RequestV1 struct {
	PriceFloor float64 `json:"price_floor"`
	PriceCeil  float64 `json:"price_ceil"`
	UserID     string  `json:"user_id"`
}

type RequestV2 struct {
	PriceLow float32 `json:"price_low"`
	User     string  `json:"user"`
	Currency string  `json:"currency"`
}

type RequestV3 string

func (r *Request) FromV1(p []byte) (err error) {
	var v1 RequestV1
	if err = json.Unmarshal(p, &v1); err != nil {
		return
	}
	r.BF, r.BC = v1.PriceFloor, v1.PriceCeil
	r.Limit = 1
	r.UID = v1.UserID
	r.Cur = "USD"
	return
}

func (r *Request) FromV2(p []byte) (err error) {
	var v2 RequestV2
	if err = json.Unmarshal(p, &v2); err != nil {
		return
	}
	r.BF = float64(v2.PriceLow)
	r.Limit = 1
	r.UID = v2.User
	r.Cur = v2.Currency
	return
}

func (r *Request) FromV3(p []byte) (err error) {
	r = &Request{}
	var v url.Values
	if v, err = url.ParseQuery(string(p)); err != nil {
		return
	}
	for k, vv := range v {
		switch k {
		case "bf":
		case "lo":
		case "pl":
			if r.BF, err = strconv.ParseFloat(vv[0], 64); err != nil {
				return
			}
		case "bc":
		case "hi":
		case "pt":
			if r.BC, err = strconv.ParseFloat(vv[0], 64); err != nil {
				return
			}
		case "limit":
			var u64 uint64
			if u64, err = strconv.ParseUint(vv[0], 10, 64); err != nil {
				return
			}
			r.Limit = uint(u64)
		case "uid":
		case "user_id":
			r.UID = vv[0]
		case "cur":
		case "currency":
			r.Cur = vv[0]
		}
	}
	return
}
