package handlers

import (
	"errors"
	"net/http"

	"github.com/identicalaffiliation/booking-service/booking/internal/dto/json"
	"github.com/labstack/echo/v4"
)

func HTTPErrorHandler() echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}

		var httpErr *json.HTTPError
		if errors.As(err, &httpErr) {
			_ = c.JSON(httpErr.Code, httpErr)
			return
		}

		_ = c.JSON(http.StatusInternalServerError, json.HTTPError{
			Message: "INTERNAL SERVER ERROR",
			Code:    http.StatusInternalServerError,
			Status:  json.INTERNAL,
		})
	}
}
