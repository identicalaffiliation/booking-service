package controller

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/identicalaffiliation/booking-service/booking/internal/domain"
	"github.com/identicalaffiliation/booking-service/booking/internal/dto/input"
	"github.com/identicalaffiliation/booking-service/booking/internal/dto/output"
	"github.com/identicalaffiliation/booking-service/booking/internal/ports"
	"github.com/labstack/echo/v4"
)

const (
	RoomIdMuxPattern = "roomId"
)

func CreateSchedule(schedule ports.SchedulesUsecase) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var in input.CreateScheduleInput
		if err := ctx.Bind(&in); err != nil {
			return output.NewBadRequest("invalid json body")
		}

		req := ctx.Request()

		roomID, err := uuid.Parse(ctx.Param(RoomIdMuxPattern))
		if err != nil {
			return output.NewBadRequest("invalid path param")
		}

		in.RoomID = roomID

		out, err := schedule.CreateSchedule(req.Context(), &in)
		if err != nil {
			if errors.Is(err, domain.ErrInvalidScheduleData) {
				return output.NewBadRequest("invalid schedule data")
			}

			if errors.Is(err, domain.ErrScheduleAlreadyExists) {
				return output.NewBadRequest("schedule already exists")
			}

			return output.NewInternal()
		}

		return ctx.JSON(http.StatusCreated, &output.CreateScheduleOutput{
			ID:        out.ID,
			RoomID:    out.RoomID,
			Day:       out.Day,
			Start:     out.Start,
			End:       out.End,
			CreatedAt: out.CreatedAt,
		})
	}
}
