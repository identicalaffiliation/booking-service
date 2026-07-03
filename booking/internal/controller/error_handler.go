package controller

import (
	"errors"
	"net/http"

	"github.com/identicalaffiliation/booking-service/booking/internal/dto/output"
	"github.com/labstack/echo/v4"
)

func HTTPErrorHandler() echo.HTTPErrorHandler {
	return func(err error, ctx echo.Context) {
		if ctx.Response().Committed {
			return
		}

		var httpErr *output.HTTPError
		if errors.As(err, &httpErr) {
			type Error struct {
				Err error `json:"error"`
			}

			err := &Error{
				Err: httpErr,
			}

			_ = ctx.JSON(httpErr.Code, err)
			return
		}

		_ = ctx.JSON(http.StatusInternalServerError, output.HTTPError{
			Message: "internal server error",
			Code:    http.StatusInternalServerError,
			Status:  output.Internal,
		})
	}
}
