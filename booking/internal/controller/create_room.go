package controller

import (
	"errors"
	"net/http"

	"github.com/identicalaffiliation/booking-service/booking/internal/domain"
	"github.com/identicalaffiliation/booking-service/booking/internal/dto/input"
	"github.com/identicalaffiliation/booking-service/booking/internal/dto/output"
	"github.com/identicalaffiliation/booking-service/booking/internal/usecase"
	"github.com/labstack/echo/v4"
)

func CreateRoom(rooms *usecase.RoomsUsecase) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var in input.CreateRoomInput
		if err := ctx.Bind(&in); err != nil {
			return output.NewBadRequest("invalid json body")
		}

		reqCtx := ctx.Request().Context()

		out, err := rooms.CreateRoom(reqCtx, &in)
		if err != nil {
			switch {
			case errors.Is(err, domain.ErrRoomAlreadyExists):
				return output.NewBadRequest("room already exists")
			case errors.Is(err, domain.ErrInvalidRoomData):
				return output.NewBadRequest("invalid room data")
			default:
				return output.NewInternal()
			}
		}

		return ctx.JSON(http.StatusCreated, &output.CreateRoomOutput{
			ID:        out.ID,
			Name:      out.Name,
			Capacity:  out.Capacity,
			CreatedAt: out.CreatedAt,
		})
	}
}
