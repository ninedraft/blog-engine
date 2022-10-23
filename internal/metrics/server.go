package metrics

import (
	"context"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Serve(ctx context.Context) {
	var server = http.Server{
		Addr:    "localhost:2112",
		Handler: promhttp.Handler(),
	}
	go func() {
		<-ctx.Done()
		_ = server.Shutdown(ctx)
	}()
	_ = server.ListenAndServe()
}
