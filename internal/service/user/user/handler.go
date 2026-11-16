package user

import (
	"go-backend-template/internal/platform/response"
	"net/http"
	"strconv"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil || id < 1 {
		return response.CodeInvalidArgument
	}
	user, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		return err
	}
	return response.Success(w, user)
}
