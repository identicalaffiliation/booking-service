package handlers

import (
	"errors"
	"net/http"

	usecase "github.com/identicalaffiliation/booking-service/booking/internal/application"
	"github.com/identicalaffiliation/booking-service/booking/internal/dto/input"
	"github.com/identicalaffiliation/booking-service/booking/internal/dto/json"
	"github.com/labstack/echo/v4"
)

func CreateRoom(roomsUsecase *usecase.RoomsUsecase) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var req json.CreateRoomRequest
		if err := ctx.Bind(&req); err != nil {
			return json.NewBadRequest("invalid json body")
		}

		in := input.NewCreateRoomInput(req.Name, req.Capacity)
		reqCtx := ctx.Request().Context()

		out, err := roomsUsecase.CreateRoom(reqCtx, in)
		if err != nil {
			switch {
			case errors.Is(err, usecase.ErrRoomAlreadyExists):
				return json.NewBadRequest("room already exists")
			case errors.Is(err, input.ErrInvalidRoomCapacity):
				return json.NewBadRequest("invalid room capacity")
			case errors.Is(err, input.ErrInvalidRoomName):
				return json.NewBadRequest("invalid room name")
			default:
				return json.NewInternal()
			}
		}

		return ctx.JSON(http.StatusCreated, &json.CreateRoomResponse{
			ID:        out.ID,
			Name:      out.Name,
			Capacity:  out.Capacity,
			CreatedAt: out.CreatedAt,
		})
	}
}
