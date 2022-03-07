package model

import (
	"encoding/json"
	"fmt"
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
	if r.BF, r.BC = v1.PriceFloor, v1.PriceCeil; r.BC <= r.BF {
		r.BC = r.BF * 3
	}
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
	r.BC = r.BF * 3
	r.Limit = 1
	r.UID = v2.User
	r.Cur = v2.Currency
	return
}

func (r *Request) FromV3(p []byte) (err error) {
	var v url.Values
	if v, err = url.ParseQuery(string(p)); err != nil {
		return
	}
	for k, vv := range v {
		switch k {
		case "bf", "lo", "pl":
			if r.BF, err = strconv.ParseFloat(vv[0], 64); err != nil {
				return
			}
		case "bc", "hi", "pt":
			if r.BC, err = strconv.ParseFloat(vv[0], 64); err != nil {
				return
			}
		case "lim", "limit":
			var u64 uint64
			if u64, err = strconv.ParseUint(vv[0], 10, 64); err != nil {
				return
			}
			r.Limit = uint(u64)
		case "uid", "user_id":
			r.UID = vv[0]
		case "cur", "currency":
			r.Cur = vv[0]
		}
	}
	if r.BC <= r.BF {
		r.BC = r.BF * 3
	}
	return
}

func (r Request) ToV1() []byte {
	v1 := RequestV1{}
	v1.PriceFloor, v1.PriceCeil = r.BF, r.BC
	v1.UserID = r.UID
	b, _ := json.Marshal(v1)
	return b
}

func (r Request) ToV2() []byte {
	v2 := RequestV2{}
	v2.PriceLow = float32(r.BF)
	v2.User = r.UID
	v2.Currency = r.Cur
	b, _ := json.Marshal(v2)
	return b
}

func (r Request) ToV3() []byte {
	return []byte(fmt.Sprintf("/v3?bf=%f&bc=%f&lim=%d&user_id=%s&cur=%s", r.BF, r.BC, r.Limit, r.UID, r.Cur))
}
