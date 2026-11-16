package user

import (
	"go-backend-template/internal/platform/response"
	"net/http"
)

var (
	CodeUserNotFound = response.NewCode(20001, "用户不存在", http.StatusNotFound)
)
