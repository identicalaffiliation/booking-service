package output

import "net/http"

const (
	BadRequest    = "BAD REQUEST"
	Internal      = "INTERNAL SERVER ERROR"
	NotFound      = "NOT FOUND"
	NotAuthorized = "UNAUTHORIZED"
	Forbidden     = "FORBIDDEN"
)

type ErrorStatus string

type HTTPError struct {
	Message string      `json:"message"`
	Code    int         `json:"code"`
	Status  ErrorStatus `json:"status"`
}

func (err *HTTPError) Error() string {
	return err.Message
}

func NewBadRequest(msg string) error {
	return &HTTPError{
		Message: msg,
		Code:    http.StatusBadRequest,
		Status:  BadRequest,
	}
}

func NewInternal() error {
	return &HTTPError{
		Message: "internal server error",
		Code:    http.StatusInternalServerError,
		Status:  Internal,
	}
}

func NewNotFound(msg string) error {
	return &HTTPError{
		Message: msg,
		Code:    http.StatusNotFound,
		Status:  NotFound,
	}
}

func NewNotAuthorized(msg string) error {
	return &HTTPError{
		Code:    http.StatusUnauthorized,
		Status:  NotAuthorized,
		Message: msg,
	}
}

func NewForbidden(msg string) error {
	return &HTTPError{
		Code:    http.StatusForbidden,
		Status:  Forbidden,
		Message: msg,
	}
}
