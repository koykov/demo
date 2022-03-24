package main

import (
	"net/http"
	"regexp"
)

type PostbackHTTP struct{}

var (
	rePB = regexp.MustCompile(`/pb/(.*)`)
)

func (h *PostbackHTTP) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// ...
}
