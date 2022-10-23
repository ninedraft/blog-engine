package router

import (
	"context"
	"strings"

	"github.com/ninedraft/gemax/gemax"
)

type Router struct {
	Routes   map[string]gemax.Handler
	Fallback gemax.Handler
}

func (router *Router) Handle(ctx context.Context, rw gemax.ResponseWriter, req gemax.IncomingRequest) {
	var requestedPath = strings.TrimSuffix(req.URL().Path, "/")

	var handle = router.Routes[requestedPath]
	if handle != nil {
		handle(ctx, rw, req)
		return
	}

	if router.Fallback != nil {
		router.Fallback(ctx, rw, req)
		return
	}

	gemax.NotFound(rw, req)
}
