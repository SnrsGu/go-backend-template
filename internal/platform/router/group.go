package router

import (
	"go-backend-template/internal/platform/middleware"
	"net/http"
	"strings"
)

type Group struct {
	router      *Router
	prefix      string
	middlewares []middleware.Middleware
}

func (g *Group) Group(prefix string) *Group {
	return &Group{
		router:      g.router,
		prefix:      joinPath(g.prefix, normalizePrefix(prefix)),
		middlewares: append([]middleware.Middleware(nil), g.middlewares...),
	}
}

func (g *Group) Use(middlewares ...middleware.Middleware) {
	g.middlewares = append(g.middlewares, middlewares...)
}

func (g *Group) Handle(method, path string, handler HandlerFunc) {
	g.router.register(method, joinPath(g.prefix, path), handler, g.middlewares)
}

func (g *Group) GET(path string, handler HandlerFunc) {
	g.Handle(http.MethodGet, path, handler)
}

func (g *Group) HEAD(path string, handler HandlerFunc) {
	g.Handle(http.MethodHead, path, handler)
}

func (g *Group) POST(path string, handler HandlerFunc) {
	g.Handle(http.MethodPost, path, handler)
}

func (g *Group) PUT(path string, handler HandlerFunc) {
	g.Handle(http.MethodPut, path, handler)
}

func (g *Group) PATCH(path string, handler HandlerFunc) {
	g.Handle(http.MethodPatch, path, handler)
}

func (g *Group) DELETE(path string, handler HandlerFunc) {
	g.Handle(http.MethodDelete, path, handler)
}

func (g *Group) CONNECT(path string, handler HandlerFunc) {
	g.Handle(http.MethodConnect, path, handler)
}

func (g *Group) OPTIONS(path string, handler HandlerFunc) {
	g.Handle(http.MethodOptions, path, handler)
}

func (g *Group) TRACE(path string, handler HandlerFunc) {
	g.Handle(http.MethodTrace, path, handler)
}

func normalizePrefix(prefix string) string {
	if prefix == "" || prefix == "/" {
		return ""
	}
	if !strings.HasPrefix(prefix, "/") {
		prefix = "/" + prefix
	}
	return strings.TrimRight(prefix, "/")
}

func joinPath(prefix string, path string) string {
	if prefix == "" || prefix == "/" {
		return "/"
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return prefix + path
}
