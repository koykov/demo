package main

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/koykov/demo/wirex/internal/dmp"
	"github.com/koykov/demo/wirex/internal/http/feed"
	"github.com/koykov/demo/wirex/internal/http/request"
	"github.com/koykov/demo/wirex/internal/models"
)

var (
	mc request.ModelContainer
)

func init() {
	mc.User = dmp.DMP{}
	mc.Source = models.SSPModel{}
}

func main() {
	router := fasthttprouter.New()
	_ = router

	handlerFeed := feed.Wire(nil, &mc)
	router.GET("/feed", handlerFeed.Process())
}
