package static

import (
	"embed"

	"github.com/ninedraft/gemax/gemax"
)

//go:embed *.gmi
var fsys embed.FS

const Route = "/info/engine"

var Handler gemax.Handler = (&gemax.FileSystem{
	FS:     fsys,
	Prefix: Route,
}).Serve
