package controller

import (
	"errors"
	"net/http"

	"github.com/identicalaffiliation/booking-service/booking/internal/dto/output"
	"github.com/labstack/echo/v4"
)

func HTTPErrorHandler() echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}

		var httpErr *output.HTTPError
		if errors.As(err, &httpErr) {
			_ = c.JSON(httpErr.Code, httpErr)
			return
		}

		_ = c.JSON(http.StatusInternalServerError, output.HTTPError{
			Message: "INTERNAL SERVER ERROR",
			Code:    http.StatusInternalServerError,
			Status:  output.INTERNAL,
		})
	}
}
