package router

import (
	"fmt"
	"net/http"
	"strings"
)

type Group struct {
	router      *Router
	prefix      string
	middlewares []Middleware
}

func (g *Group) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	g.router.ServeHTTP(writer, request)
}

func (g *Group) Group(prefix string) *Group {
	validPath(prefix)
	middlewares := append([]Middleware(nil), g.middlewares...)
	return &Group{
		router:      g.router,
		prefix:      joinPath(g.prefix, prefix),
		middlewares: middlewares,
	}
}

func (g *Group) Use(middlewares ...Middleware) {
	for _, middleware := range middlewares {
		if middleware == nil {
			panic("使用了空中间件")
		}
		g.middlewares = append(g.middlewares, middleware)
	}
}

func (g *Group) Handle(method, path string, handler http.Handler) {
	if handler == nil {
		panic("路由Handler为空")
	}
	if method == "" || strings.ContainsAny(method, " \t\r\n") {
		panic(fmt.Sprintf("请求方式不合法: %q", method))
	}
	validPath(path)

	fullPath := joinPath(g.prefix, path)
	if strings.HasSuffix(fullPath, "/") {
		fullPath += "{$}"
	}

	wrappedHandler := chain(handler, g.middlewares)

	g.router.mux.Handle(method+" "+fullPath, wrappedHandler)
}

func (g *Group) HandleFunc(method, path string, handler http.HandlerFunc) {
	if handler == nil {
		panic("路由HandlerFunc为空")
	}
	g.Handle(method, path, handler)
}

func (g *Group) GET(path string, handler http.HandlerFunc) {
	g.HandleFunc(http.MethodGet, path, handler)
}

func (g *Group) HEAD(path string, handler http.HandlerFunc) {
	g.HandleFunc(http.MethodHead, path, handler)
}

func (g *Group) POST(path string, handler http.HandlerFunc) {
	g.HandleFunc(http.MethodPost, path, handler)
}

func (g *Group) PUT(path string, handler http.HandlerFunc) {
	g.HandleFunc(http.MethodPut, path, handler)
}

func (g *Group) PATCH(path string, handler http.HandlerFunc) {
	g.HandleFunc(http.MethodPatch, path, handler)
}

func (g *Group) DELETE(path string, handler http.HandlerFunc) {
	g.HandleFunc(http.MethodDelete, path, handler)
}

func (g *Group) CONNECT(path string, handler http.HandlerFunc) {
	g.HandleFunc(http.MethodConnect, path, handler)
}

func (g *Group) OPTIONS(path string, handler http.HandlerFunc) {
	g.HandleFunc(http.MethodOptions, path, handler)
}

func (g *Group) TRACE(path string, handler http.HandlerFunc) {
	g.HandleFunc(http.MethodTrace, path, handler)
}

func joinPath(prefix, path string) string {
	if path == "" {
		if prefix == "" {
			return "/"
		}
		return prefix
	}
	if path == "/" {
		if prefix == "" {
			return "/"
		}
		return prefix + "/"
	}
	return prefix + path
}

func validPath(path string) {
	if path == "" {
		return
	}
	if !strings.HasPrefix(path, "/") {
		panic(fmt.Sprintf("路由路径必须以 / 开头: %q", path))
	}
	if strings.Contains(path, "//") {
		panic(fmt.Sprintf("路由路径不能包含连续斜杠: %q", path))
	}

	if strings.ContainsAny(path, " \t\r\n?#") {
		panic(fmt.Sprintf("路由路径包含非法字符: %q", path))
	}
}
