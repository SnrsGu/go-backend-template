package user

import "go-backend-template/internal/platform/router"

func RegisterRoutes(group *router.Group, handler *Handler) {
	user := group.Group("/user")
	user.POST("", handler.Create)
	user.GET("", handler.List)
	user.GET("/{id}", handler.GetByID)
	user.PATCH("/{id}", handler.UpdateUsername)
	user.DELETE("/{id}", handler.Delete)
}
