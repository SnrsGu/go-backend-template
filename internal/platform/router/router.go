package router

import (
	"go-backend-template/internal/platform/response"
	"log/slog"
	"net/http"
)

type Router struct {
	mux    *http.ServeMux
	logger *slog.Logger
}

func New(logger *slog.Logger) *Group {
	if logger == nil {
		panic("路由日志实例不可为空")
	}
	return &Group{
		router: &Router{
			mux:    http.NewServeMux(),
			logger: logger,
		},
	}
}

func (r *Router) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	fallback, pattern := r.mux.Handler(request)
	if pattern != "" {
		r.mux.ServeHTTP(writer, request)
		return
	}
	probe := &fallbackWriter{
		header: make(http.Header),
	}
	fallback.ServeHTTP(probe, request)
	var code response.Code
	switch probe.status {
	case http.StatusMethodNotAllowed:
		code = response.CodeMethodNotAllowed
		if allow := probe.header.Values("Allow"); len(allow) > 0 {
			for _, a := range allow {
				writer.Header().Add("Allow", a)
			}
		}
	default:
		code = response.CodeMethodNotAllowed
	}
	if err := response.Fail(writer, code); err != nil {
		r.logger.Error("错误响应失败", "error", err)
	}
}

type fallbackWriter struct {
	header http.Header
	status int
}

func (f *fallbackWriter) Header() http.Header {
	return f.header
}

func (f *fallbackWriter) WriteHeader(status int) {
	if f.status == 0 {
		f.status = status
	}
}

func (f *fallbackWriter) Write(b []byte) (int, error) {
	if f.status == 0 {
		f.status = http.StatusOK
	}
	return len(b), nil
}
