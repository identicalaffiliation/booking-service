package controller

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/identicalaffiliation/booking-service/booking/internal/domain"
	"github.com/identicalaffiliation/booking-service/booking/internal/dto/output"
	"github.com/identicalaffiliation/booking-service/booking/internal/ports"
	"github.com/labstack/echo/v4"
)

const (
	bookingIdMuxPattern = "bookingId"
)

func CancelBooking(booking ports.BookingsUsecase) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		idStr := ctx.Param(bookingIdMuxPattern)

		id, err := uuid.Parse(idStr)
		if err != nil {
			return output.NewBadRequest("invalid booking data")
		}

		reqCtx := ctx.Request().Context()
		if err := booking.Cancel(reqCtx, id); err != nil {
			if errors.Is(err, domain.ErrInvalidBookingData) {
				return output.NewBadRequest("invalid booking data")
			}

			return output.NewInternal()
		}

		return ctx.NoContent(http.StatusOK)
	}
}
