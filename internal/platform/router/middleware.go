package router

import "slices"

import "net/http"

type Middleware func(http.Handler) http.Handler

func chain(handler http.Handler, middlewares []Middleware) http.Handler {
	for _, middleware := range slices.Backward(middlewares) {
		handler = middleware(handler)
	}
	return handler
}
