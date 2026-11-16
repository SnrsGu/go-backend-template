package user

import "go-backend-template/internal/platform/router"

func RegisterRoutes(group *router.Group, handler *Handler) {
	user := group.Group("/api/user")
	user.GET("/{id}", handler.GetByID)
}
