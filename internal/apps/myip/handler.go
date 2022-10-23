package myip

import (
	"context"
	"fmt"

	"github.com/ninedraft/gemax/gemax"
)

func Handle(_ context.Context, rw gemax.ResponseWriter, req gemax.IncomingRequest) {
	fmt.Fprintf(rw, "%s\r\n", req.RemoteAddr())
}
