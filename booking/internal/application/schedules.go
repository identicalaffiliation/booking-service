package application

import (
		"context"
		"errors"
		"fmt"

		"github.com/google/uuid"
		"github.com/identicalaffiliation/booking-service/booking/internal/adapters/storage/psql"
		"github.com/identicalaffiliation/booking-service/booking/internal/domain"
		"github.com/identicalaffiliation/booking-service/booking/internal/dto/input"
		"github.com/identicalaffiliation/booking-service/booking/internal/dto/output"
		"github.com/identicalaffiliation/booking-service/booking/internal/ports"
)

const ScheduleLayer = "schedule service layer"

type SchedulesUsecase struct {
		schRepo ports.SchedulesRepository
		log     ports.Logger
}

func NewSchedulesUsecase(schRepo ports.SchedulesRepository, log ports.Logger) *SchedulesUsecase {
		return &SchedulesUsecase{
				schRepo: schRepo,
				log:     log,
		}
}

func (u *SchedulesUsecase) CreateSchedule(ctx context.Context, in *input.CreateScheduleInput) (*output.CreateScheduleOutput, error) {
		if in.RoomID == uuid.Nil {
				return nil, ErrInvalidRoomId
		}

		parsedStartTime, err := domain.ParseTimeDuration(in.Start)
		if err != nil {
				return nil, ErrInternalScheduleUsecase
		}

		parsedEndTime, err := domain.ParseTimeDuration(in.End)
		if err != nil {
				return nil, ErrInternalScheduleUsecase
		}

		parsedDay, err := domain.ParseTimeDate(in.Day)
		if err != nil {
				return nil, ErrInternalScheduleUsecase
		}

		if !u.ValidateStartAndEndInterval(parsedStartTime, parsedEndTime) {
				return nil, ErrInvalidTimeInterval
		}

		schedule := domain.NewSchedule(
				in.RoomID,
				parsedDay,
				parsedStartTime,
				parsedEndTime,
		)

		created, err := u.schRepo.CreateSchedule(ctx, schedule)
		if err != nil {
				if errors.Is(err, psql.ErrScheduleAlreadyExists) {
						return nil, ErrScheduleAlreadyExists
				}

				u.log.Error("failed to create schedule", "layer", ScheduleLayer, "error", err)
				return nil, ErrInternalScheduleUsecase
		}

		return output.NewCreateScheduleOutput(
				created.ID,
				created.RoomID,
				created.Day.Format("2006-01-02"),
				u.minutesToTime(created.StartWorkTime),
				u.minutesToTime(created.EndWorkTime),
				created.CreatedAt.UTC(),
		), nil
}

func (u *SchedulesUsecase) minutesToTime(minutes int) string {
		h := minutes / 60
		m := minutes % 60

		return fmt.Sprintf("%02d:%02d", h, m)
}

func (u *SchedulesUsecase) ValidateStartAndEndInterval(start, end int) bool {
		if start >= end {
				return false
		}

		if start%60 != 0 || end%60 != 0 {
				return false
		}

		return true
}
