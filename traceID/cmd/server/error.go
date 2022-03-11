package main

import "errors"

var (
	ErrBidfloorFail = errors.New("bid less than bidfloor")
	ErrBidceilFail  = errors.New("bid overflows bidceil")
)
