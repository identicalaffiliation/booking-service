package controller

import (
	"net/http"

	"github.com/identicalaffiliation/booking-service/booking/internal/dto/output"
	"github.com/identicalaffiliation/booking-service/booking/internal/usecase"
	"github.com/labstack/echo/v4"
)

func GetRooms(rooms *usecase.RoomsUsecase) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		reqCtx := ctx.Request().Context()
		out, err := rooms.GetRooms(reqCtx)
		if err != nil {
			return output.NewInternal()
		}

		return ctx.JSON(http.StatusOK, out)
	}
}
