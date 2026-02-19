package apperrors

import "net/http"

type AppError struct {
	Code    string
	Message string
	Status  int
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

func (e *AppError) Unwrap() error { return e.Err }

func ErrUnauthorized(msg string) *AppError {
	return &AppError{Code: "unauthorized", Message: msg, Status: http.StatusUnauthorized}
}

func ErrForbidden(msg string) *AppError {
	return &AppError{Code: "forbidden", Message: msg, Status: http.StatusForbidden}
}

func ErrNotFound(msg string) *AppError {
	return &AppError{Code: "not_found", Message: msg, Status: http.StatusNotFound}
}

func ErrInvalid(msg string) *AppError {
	return &AppError{Code: "invalid", Message: msg, Status: http.StatusBadRequest}
}

func ErrInternal(msg string) *AppError {
	return &AppError{Code: "internal", Message: msg, Status: http.StatusInternalServerError}
}

func ErrConflict(msg string) *AppError {
	return &AppError{Code: "conflict", Message: msg, Status: http.StatusConflict}
}