package controller

import (
	"errors"
	"net/http"

	"github.com/identicalaffiliation/booking-service/booking/internal/domain"
	"github.com/identicalaffiliation/booking-service/booking/internal/dto/input"
	"github.com/identicalaffiliation/booking-service/booking/internal/dto/output"
	"github.com/identicalaffiliation/booking-service/booking/internal/ports"
	"github.com/labstack/echo/v4"
)

func CreateBooking(booking ports.BookingsUsecase) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var in input.CreateBookingInput
		if err := ctx.Bind(&in); err != nil {
			return output.NewBadRequest("invalid json body")
		}

		reqCtx := ctx.Request().Context()

		out, err := booking.Create(reqCtx, &in)
		if err != nil {
			switch {
			case errors.Is(err, domain.ErrInvalidBookingData):
				return output.NewBadRequest("invalid booking data")
			case errors.Is(err, domain.ErrInvalidUserData):
				return output.NewBadRequest("invalid user data")
			case errors.Is(err, domain.ErrSlotAlreadyBooked):
				return output.NewBadRequest("slot already booked")
			default:
				return output.NewInternal()
			}
		}
		
		return ctx.JSON(http.StatusCreated, out)
	}
}
