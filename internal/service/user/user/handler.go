package user

import (
	"encoding/json"
	"errors"
	"go-backend-template/internal/platform/response"
	"io"
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

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) error {
	var req CreateUserRequest

	if err := decodeJSON(w, r, &req); err != nil {
		return err
	}

	user, err := h.service.Create(r.Context(), req)
	if err != nil {
		return err
	}

	return response.Success(w, user)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) error {
	query := r.URL.Query()

	page := 1
	pageSize := 20

	if value := query.Get("page"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil {
			return response.CodeInvalidArgument
		}
		page = parsed
	}

	if value := query.Get("page_size"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil {
			return response.CodeInvalidArgument
		}
		pageSize = parsed
	}

	// 限制分页范围，避免不受控的大量数据查询。
	if page < 1 || page > 1000000 ||
		pageSize < 1 || pageSize > 100 {
		return response.CodeInvalidArgument
	}

	result, err := h.service.List(r.Context(), page, pageSize)
	if err != nil {
		return err
	}

	return response.Success(w, result)
}

func (h *Handler) UpdateUsername(w http.ResponseWriter, r *http.Request) error {
	id, err := parseUserID(r)
	if err != nil {
		return err
	}

	var req UpdateUserRequest

	if err := decodeJSON(w, r, &req); err != nil {
		return err
	}

	if err := h.service.UpdateUsername(
		r.Context(),
		id,
		req,
	); err != nil {
		return err
	}

	return response.Success(w, nil)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) error {
	id, err := parseUserID(r)
	if err != nil {
		return err
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		return err
	}

	return response.Success(w, nil)
}

func parseUserID(r *http.Request) (uint64, error) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil || id == 0 {
		return 0, response.CodeInvalidArgument
	}

	return id, nil
}

func decodeJSON(w http.ResponseWriter, r *http.Request, target any) error {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(target); err != nil {
		return response.CodeInvalidArgument
	}

	// 不允许一个请求体包含多个 JSON 对象。
	var extra any

	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		return response.CodeInvalidArgument
	}

	return nil
}
