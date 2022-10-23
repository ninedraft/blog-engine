package middlewares

import "github.com/ninedraft/gemax/gemax"

type Middleware func(handle gemax.Handler) gemax.Handler

func With(handle gemax.Handler, mws ...Middleware) gemax.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		handle = mws[i](handle)
	}
	return handle
}
