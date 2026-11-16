package gateway

import "go-backend-template/internal/platform/router"

func RegisterRouters(appRouter *router.Group) {
	appRouter.GET("/health", HealthHandler)
}
