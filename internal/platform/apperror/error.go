package apperror

import (
	"errors"
	"fmt"
)

type AppError struct {
	Code    Code
	Message string
	Cause   error
}

func New(code Code, msg string) *AppError {
	return &AppError{
		Code:    code,
		Message: msg,
	}
}

func Wrap(code Code, msg string, cause error) *AppError {
	return &AppError{
		Code:    code,
		Message: msg,
		Cause:   cause,
	}
}

func FromError(err error) *AppError {
	var e *AppError
	if errors.As(err, &e) && e != nil {
		normalizedCode := e.Code.Normalize()
		if normalizedCode != e.Code {
			return Wrap(CodeInternal, "", err)
		}
		return e
	}
	return Wrap(CodeInternal, "", err)
}

func (e *AppError) Unwrap() error {
	return e.Cause
}

func (e *AppError) Error() string {
	msg := e.PublicMessage()
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", msg, e.Cause)
	}
	return msg
}

func (e *AppError) PublicMessage() string {
	if e.Message != "" {
		return e.Message
	}
	return e.Code.Message()
}
