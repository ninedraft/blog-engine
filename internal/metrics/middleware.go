package metrics

import (
	"context"

	"github.com/ninedraft/gemax/gemax"
)

func With(handle gemax.Handler) gemax.Handler {
	return func(ctx context.Context, rw gemax.ResponseWriter, req gemax.IncomingRequest) {
		pageCounter.WithLabelValues(req.URL().Path).Inc()
		handle(ctx, rw, req)
	}
}
