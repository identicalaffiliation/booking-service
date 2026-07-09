package controller

import (
	"errors"
	"net/http"

	"github.com/identicalaffiliation/booking-service/booking/internal/domain"
	"github.com/identicalaffiliation/booking-service/booking/internal/dto/output"
	"github.com/identicalaffiliation/booking-service/booking/internal/ports"
	"github.com/labstack/echo/v4"
)

const (
	BookingIdMuxPattern = "bookingId"
)

func GetBooking(booking ports.BookingsUsecase) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		bookingIdStrFromParams := ctx.Param(BookingIdMuxPattern)

		reqCtx := ctx.Request().Context()
		out, err := booking.GetMyBooking(reqCtx, bookingIdStrFromParams)
		if err != nil {
			switch {
			case errors.Is(err, domain.ErrInvalidBookingData):
				return output.NewBadRequest("invalid booking data")
			case errors.Is(err, domain.ErrBookingNotFound):
				return output.NewNotFound("booking not found")
			default:
				return output.NewInternal()
			}
		}

		return ctx.JSON(http.StatusOK, out)
	}
}
