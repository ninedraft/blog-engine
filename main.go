package main

import (
	"archive/zip"
	"context"
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"runtime/debug"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/ninedraft/blog-engine/internal/apps/myip"
	"github.com/ninedraft/blog-engine/internal/metrics"
	"github.com/ninedraft/blog-engine/internal/middlewares"
	"github.com/ninedraft/blog-engine/internal/render"
	"github.com/ninedraft/blog-engine/internal/router"
	"github.com/ninedraft/blog-engine/static"

	"github.com/ninedraft/gemax/gemax"
)

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

	var contentSource = "content"
	flag.StringVar(&contentSource, "content", contentSource, "content source. Can be dir, a file or a zip archive with files")

	var enableRender = false
	flag.BoolVar(&enableRender, "render", enableRender, "enable go template content rendering")

	var renderGlobs = []string{"*.gmi"}
	flag.Func("render-glob",
		"allowlist for renderer. NOOP if -render is not set. Default values: "+strings.Join(renderGlobs, ", "),
		func(value string) error {
			renderGlobs = append(renderGlobs, value)

			return nil
		})

	flag.Parse()

	var version = "dev"
	if info, _ := debug.ReadBuildInfo(); info != nil {
		version = info.Main.Version
	}

	if user != "" {
		user = "/~" + user
	}

	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill)
	defer cancel()

	go metrics.Serve(ctx)

	log.Printf("building content fs...")
	var content, closeFs, errContent = buildContent(contentSource)
	if closeFs != nil {
		defer closeFs()
	}

	if errContent != nil {
		panic("reading content source" + errContent.Error())
	}

	log.Printf("building content fs - ok")

	if enableRender {
		log.Printf("rendering...")
		var renderFS, errRender = render.New(content, renderGlobs...)
		if errRender != nil {
			globs := strings.Join(renderGlobs, ", ")
			panic("preparing renderer for files " + globs + ": " + errRender.Error())
		}
		content = renderFS
		log.Printf("rendering - ok")
	}

	var handleContent = (&gemax.FileSystem{
		Prefix: "",
		FS:     content,
		Logf:   log.Printf,
	}).Serve

	var routes = router.Router{
		Group: user,
		Routes: map[string]gemax.Handler{
			"/myip": myip.Handle,
			"/version": gemax.ServeContent(gemax.MIMEGemtext,
				[]byte(fmt.Sprintf("# Engine version\n\n%s", version))),
			static.Route: static.Handler,
		},
		Fallback: handleContent,
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

	log.Printf("serving...")
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

func buildContent(filename string) (_ fs.FS, closeFs func(), _ error) {
	var stat, errStat = os.Stat(filename)
	if errStat != nil {
		return nil, nil, fmt.Errorf("checking file stats: %w", errStat)
	}

	if stat.IsDir() {
		return os.DirFS(filename), func() {}, nil
	}

	if filepath.Ext(filename) != ".zip" {
		return nil, nil, fmt.Errorf("%w: only zip files and dirs are supported", errors.ErrUnsupported)
	}

	const flags = os.O_RDONLY | syscall.O_NOFOLLOW
	const perm = 0666
	var file, errFile = os.OpenFile(filename, flags, perm)
	if errFile != nil {
		return nil, nil, fmt.Errorf("opening zip file: %w", errFile)
	}
	closeFs = sync.OnceFunc(func() {
		_ = file.Close()
	})

	var archive, errArchive = zip.NewReader(file, stat.Size())
	if errArchive != nil {
		return nil, closeFs, fmt.Errorf("reading zip file: %w", errArchive)
	}

	return archive, closeFs, nil
}
