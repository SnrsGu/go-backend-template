package handler

import (
	"go-backend-template/internal/platform/response"
	"net/http"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) error {
	return response.Success(w, map[string]string{"status": "ok"})
}
