package controller

import (
	"errors"
	"net/http"

	"github.com/identicalaffiliation/booking-service/booking/internal/domain"
	"github.com/identicalaffiliation/booking-service/booking/internal/dto/output"
	"github.com/identicalaffiliation/booking-service/booking/internal/ports"
	"github.com/labstack/echo/v4"
)

func GetRoom(rooms ports.RoomsUsecase) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		id := ctx.Param(RoomIdMuxPattern)

		reqCtx := ctx.Request().Context()

		out, err := rooms.GetRoom(reqCtx, id)
		if err != nil {
			switch {
			case errors.Is(err, domain.ErrInvalidRoomData):
				return output.NewBadRequest("invalid room data")
			case errors.Is(err, domain.ErrRoomNotFound):
				return output.NewNotFound("room not found")
			default:
				return output.NewInternal()
			}
		}

		return ctx.JSON(http.StatusOK, out)
	}
}
