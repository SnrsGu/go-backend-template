package router

import "net/http"

type Router struct {
	mux *http.ServeMux
}

func New() *Group {
	return &Group{
		router: &Router{
			mux: http.NewServeMux(),
		},
	}
}

func (r *Router) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	r.mux.ServeHTTP(writer, request)
}
