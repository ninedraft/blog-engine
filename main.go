package main

import (
	"context"
	"crypto/tls"
	"embed"
	"errors"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/ninedraft/blog-engine/internal/apps/myip"
	"github.com/ninedraft/blog-engine/internal/metrics"
	"github.com/ninedraft/blog-engine/internal/middlewares"
	"github.com/ninedraft/blog-engine/internal/router"

	"github.com/ninedraft/gemax/gemax"
)

//go:embed content
var content embed.FS

func main() {
	var addr = "localhost:1965"
	flag.StringVar(&addr, "addr", addr, "optional address to serve")

	var caCert = "cert.pem"
	flag.StringVar(&caCert, "ca-cert", caCert, "certificate file")

	var caKey = "key.pem"
	flag.StringVar(&caKey, "ca-key", caKey, "certificate key")

	var host string
	flag.StringVar(&host, "host", host, "optional host")

	var user = os.Getenv("USER")
	flag.StringVar(&user, "user", user, "user prefix to paths: ~$USER")

	flag.Parse()

	if user != "" {
		user = "/~" + user
	}

	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill)
	defer cancel()

	go metrics.Serve(ctx)

	var handleContent = (&gemax.FileSystem{
		Prefix: "content",
		FS:     content,
		Logf:   log.Printf,
	}).Serve

	var routes = router.Router{
		Group: user,
		Routes: map[string]gemax.Handler{
			"/myip": myip.Handle,
		},
		Fallback: func(ctx context.Context, rw gemax.ResponseWriter, req gemax.IncomingRequest) {
			handleContent(ctx, rw, req)
		},
	}

	var server = gemax.Server{
		Addr: addr,
		Handler: middlewares.With(routes.Handle,
			metrics.With,
		),
		Logf: log.Printf,
	}

	if host != "" {
		server.Hosts = append(server.Hosts, host)
	}

	go func() {
		<-ctx.Done()
		server.Stop()
		time.Sleep(time.Second)
		panic("watchdog timeout")
	}()

	var cert, errCert = tls.LoadX509KeyPair(caCert, caKey)
	if errCert != nil {
		panic("loading cert: " + errCert.Error())
	}

	var errServe = server.ListenAndServe(ctx, &tls.Config{
		Certificates: []tls.Certificate{cert},
	})
	switch {
	case errors.Is(errServe, net.ErrClosed):
		return
	case errServe != nil:
		panic("serving: " + errServe.Error())
	}
}
