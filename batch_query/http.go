package main

import "net/http"

type BQHTTP struct{}

func (h *BQHTTP) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.URL.Path == "/api/v1/ping":
		// pong
	}
}
