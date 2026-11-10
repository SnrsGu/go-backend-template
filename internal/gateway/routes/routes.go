package routes

import (
	"go-backend-template/internal/gateway/handler"
	"go-backend-template/internal/platform/router"
)

func Register(router *router.Router) {
	healthHandler := handler.NewHealthHandler()
	router.GET("/health", healthHandler.Check)

	apiV1 := router.Group("/api/v1")
	apiV1.GET("/health", healthHandler.Check)
}
