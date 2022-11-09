package router

import (
	"context"
	"strings"

	"github.com/ninedraft/gemax/gemax"
)

type Router struct {
	Group    string
	Routes   map[string]gemax.Handler
	Fallback gemax.Handler
}

func (router *Router) Handle(ctx context.Context, rw gemax.ResponseWriter, req gemax.IncomingRequest) {
	var requestedPath = strings.TrimSuffix(req.URL().Path, "/")

	if !strings.HasPrefix(requestedPath, router.Group) {
		gemax.NotFound(rw, req)
		return
	}
	requestedPath = strings.TrimPrefix(requestedPath, router.Group)

	var handle = router.Routes[requestedPath]
	if handle != nil {
		handle(ctx, rw, req)
		return
	}

	if router.Fallback != nil {
		req.URL().Path = requestedPath
		router.Fallback(ctx, rw, req)
		return
	}

	gemax.NotFound(rw, req)
}
