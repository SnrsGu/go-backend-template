package response

import (
	"encoding/json/v2"
	"go-backend-template/internal/platform/apperror"
	"net/http"
)

type Result struct {
	Code    apperror.Code `json:"code"`
	Message string        `json:"message"`
	Data    any           `json:"data"`
}

func JSON(w http.ResponseWriter, httpStatus int, result Result) error {
	respBody, err := json.Marshal(result)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	_, err = w.Write(append(respBody, '\n'))
	return err
}

func Success(w http.ResponseWriter, data any) error {
	return JSON(w, http.StatusOK, Result{
		Code:    apperror.CodeSuccess,
		Message: apperror.CodeSuccess.Message(),
		Data:    data,
	})
}

func Failure(w http.ResponseWriter, err error) error {
	appError := apperror.FromError(err)
	return JSON(w, appError.Code.HTTPStatus(), Result{
		Code:    appError.Code,
		Message: appError.PublicMessage(),
		Data:    nil,
	})
}
