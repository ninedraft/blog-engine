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

	flag.Parse()

	var routes = router.Router{
		Routes: map[string]gemax.Handler{
			"/myip": myip.Handle,
		},
		Fallback: (&gemax.FileSystem{
			Prefix: "content",
			FS:     content,
			Logf:   log.Printf,
		}).Serve,
	}

	var server = gemax.Server{
		Addr:    addr,
		Handler: routes.Handle,
		Logf:    log.Printf,
	}

	if host != "" {
		server.Hosts = append(server.Hosts, host)
	}

	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill)
	defer cancel()

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
