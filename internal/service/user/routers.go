package user

import (
	"go-backend-template/internal/platform/router"
	usermodule "go-backend-template/internal/service/user/user"
)

func RegisterRoutes(group *router.Group, handler *usermodule.Handler) {
	api := group.Group("/api")
	usermodule.RegisterRoutes(api, handler)
}
