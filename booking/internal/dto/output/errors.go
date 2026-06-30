package output

import "net/http"

const (
	BadRequest = "BAD_REQUEST"
	INTERNAL   = "INTERNAL_SERVER_ERROR"
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
		Code:   http.StatusInternalServerError,
		Status: INTERNAL,
	}
}
