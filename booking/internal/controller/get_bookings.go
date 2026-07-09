package controller

import (
	"errors"
	"net/http"

	"github.com/identicalaffiliation/booking-service/booking/internal/domain"
	"github.com/identicalaffiliation/booking-service/booking/internal/dto/output"
	"github.com/identicalaffiliation/booking-service/booking/internal/ports"
	"github.com/labstack/echo/v4"
)

func GetBookings(booking ports.BookingsUsecase) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		reqCtx := ctx.Request().Context()
		out, err := booking.GetMyBookings(reqCtx)
		if err != nil {
			switch {
			case errors.Is(err, domain.ErrInvalidUserData):
				return output.NewBadRequest("invalid user data")
			default:
				return output.NewInternal()
			}
		}
		
		return ctx.JSON(http.StatusOK, out)
	}
}
