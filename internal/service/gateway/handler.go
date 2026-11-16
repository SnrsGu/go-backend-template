package gateway

import (
	"go-backend-template/internal/platform/response"
	"net/http"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) error {
	return response.Success(w, map[string]string{"status": "ok"})
}
