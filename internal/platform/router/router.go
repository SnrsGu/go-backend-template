package router

import (
	"go-backend-template/internal/platform/apperror"
	"go-backend-template/internal/platform/response"
	"log/slog"
	"net/http"
)

type HandlerFunc func(http.ResponseWriter, *http.Request) error

type Router struct {
	mux    *http.ServeMux
	logger *slog.Logger
}

func New(logger *slog.Logger) *Router {
	return &Router{
		mux:    http.NewServeMux(),
		logger: logger,
	}
}

func (r *Router) Handle(method string, path string, handler HandlerFunc) {
	pattern := method + " " + path
	r.mux.HandleFunc(pattern, func(w http.ResponseWriter, req *http.Request) {
		if err := handler(w, req); err != nil {
			r.handleError(w, req, err)
		}
	})
}

func (r *Router) Handler() http.Handler {
	return r.mux
}

func (r *Router) handleError(w http.ResponseWriter, req *http.Request, err error) {
	appError := apperror.FromError(err)
	logMethod := r.logger.Warn
	if appError.Code.HTTPStatus() >= http.StatusInternalServerError {
		logMethod = r.logger.Error
	}
	logMethod(
		"HTTP handler error",
		"method", req.Method,
		"path", req.URL.Path,
		"code", appError.Code,
		"error", err,
	)

	if writeErr := response.Failure(w, appError); writeErr != nil {
		r.logger.Error("failed to write error response", "error", writeErr)
	}
}

func (r *Router) GET(path string, handler HandlerFunc) {
	r.Handle(http.MethodGet, path, handler)
}

func (r *Router) HEAD(path string, handler HandlerFunc) {
	r.Handle(http.MethodHead, path, handler)
}

func (r *Router) POST(path string, handler HandlerFunc) {
	r.Handle(http.MethodPost, path, handler)
}

func (r *Router) PUT(path string, handler HandlerFunc) {
	r.Handle(http.MethodPut, path, handler)
}

func (r *Router) PATCH(path string, handler HandlerFunc) {
	r.Handle(http.MethodPatch, path, handler)
}

func (r *Router) DELETE(path string, handler HandlerFunc) {
	r.Handle(http.MethodDelete, path, handler)
}

func (r *Router) CONNECT(path string, handler HandlerFunc) {
	r.Handle(http.MethodConnect, path, handler)
}

func (r *Router) OPTIONS(path string, handler HandlerFunc) {
	r.Handle(http.MethodOptions, path, handler)
}

func (r *Router) TRACE(path string, handler HandlerFunc) {
	r.Handle(http.MethodTrace, path, handler)
}
