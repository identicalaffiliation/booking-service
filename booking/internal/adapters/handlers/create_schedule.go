package handlers

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/identicalaffiliation/booking-service/booking/internal/application"
	"github.com/identicalaffiliation/booking-service/booking/internal/dto/input"
	"github.com/identicalaffiliation/booking-service/booking/internal/dto/json"
	"github.com/labstack/echo/v4"
)

const (
	ROOM_ID_MUX_PATTERN = "roomId"
)

func CreateSchedule(schedule *application.SchedulesUsecase) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var r json.CreateScheduleRequest
		if err := ctx.Bind(&r); err != nil {
			return json.NewBadRequest("invalid json body")
		}

		req := ctx.Request()

		roomID, err := uuid.Parse(ctx.Param(ROOM_ID_MUX_PATTERN))
		if err != nil {
			return json.NewBadRequest("invalid room id")
		}

		in := input.NewCreateScheduleInput(roomID, r.Day, r.Start, r.End)

		out, err := schedule.CreateSchedule(req.Context(), in)
		if err != nil {
			if errors.Is(err, application.ErrInvalidRoomId) {
				return json.NewBadRequest("invalid room id")
			}

			if errors.Is(err, application.ErrInvalidTimeInterval) {
				return json.NewBadRequest("invalid time interval")
			}

			if errors.Is(err, application.ErrScheduleAlreadyExists) {
				return json.NewBadRequest("schedule already exists")
			}

			return json.NewInternal()
		}

		return ctx.JSON(http.StatusCreated, &json.CreateScheduleResponse{
			ID:        out.ID,
			RoomID:    out.RoomID,
			Day:       out.Day,
			Start:     out.Start,
			End:       out.End,
			CreatedAt: out.CreatedAt,
		})
	}
}
