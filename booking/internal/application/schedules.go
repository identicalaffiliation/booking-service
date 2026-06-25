package application

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/identicalaffiliation/booking-service/booking/internal/adapters/storage/psql"
	"github.com/identicalaffiliation/booking-service/booking/internal/domain"
	"github.com/identicalaffiliation/booking-service/booking/internal/dto/input"
	"github.com/identicalaffiliation/booking-service/booking/internal/dto/output"
	"github.com/identicalaffiliation/booking-service/booking/internal/ports"
)

var (
	ErrScheduleAlreadyExists   = errors.New("schedule already exists")
	ErrInvalidRoomId           = errors.New("invalid room id")
	ErrInvalidTimeInterval     = errors.New("invalid start and end time interval")
	ErrInternalScheduleUsecase = errors.New("invalid server error")
)

type SchedulesUsecase struct {
	schRepo ports.SchedulesRepository
}

func NewSchedulesUsecase(schRepo ports.SchedulesRepository) *SchedulesUsecase {
	return &SchedulesUsecase{
		schRepo: schRepo,
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

		return nil, ErrInternalScheduleUsecase
	}

	return output.NewCreateScheduleOutput(
		created.ID,
		created.RoomID,
		created.Day.UTC().GoString(),
		created.StartWorkTime%60,
		created.EndWorkTime%60,
		created.CreatedAt.UTC(),
	), nil
}

func (u *SchedulesUsecase) ValidateStartAndEndInterval(start, end int) bool {
	if start >= end {
		return false
	}

	if start%60 != 0 {
		return false
	}

	if end%60 != 0 {
		return false
	}

	return true
}
