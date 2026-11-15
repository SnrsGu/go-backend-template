package router

import (
	"log/slog"
	"net/http"
)

type Router struct {
	mux    *http.ServeMux
	logger *slog.Logger
}

func New(logger *slog.Logger) *Group {
	return &Group{
		router: &Router{
			mux:    http.NewServeMux(),
			logger: logger,
		},
	}
}

func (r *Router) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	r.mux.ServeHTTP(writer, request)
}
