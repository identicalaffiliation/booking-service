package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/identicalaffiliation/booking-service/booking/internal/domain"
	"github.com/identicalaffiliation/booking-service/booking/internal/dto/input"
	"github.com/identicalaffiliation/booking-service/booking/internal/dto/output"
	"github.com/identicalaffiliation/booking-service/booking/internal/ports"
)

const ScheduleLayer = "schedule service layer"

type SchedulesUsecase struct {
	schRepo ports.SchedulesRepository
	log     ports.Logger
	slots   *SlotsUsecase
}

func NewSchedulesUsecase(schRepo ports.SchedulesRepository, log ports.Logger, slots *SlotsUsecase) *SchedulesUsecase {
	return &SchedulesUsecase{
		schRepo: schRepo,
		log:     log,
		slots:   slots,
	}
}

func (u *SchedulesUsecase) CreateSchedule(ctx context.Context,
	in *input.CreateScheduleInput,
) (*output.CreateScheduleOutput, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	parsedStartTime, err := domain.ParseTimeDuration(in.Start)
	if err != nil {
		return nil, domain.ErrInvalidScheduleData
	}

	parsedEndTime, err := domain.ParseTimeDuration(in.End)
	if err != nil {
		return nil, domain.ErrInvalidScheduleData
	}

	parsedDay, err := domain.ParseTimeDate(in.Day)
	if err != nil {
		return nil, domain.ErrInvalidScheduleData
	}

	schedule := domain.NewSchedule(
		in.RoomID,
		parsedDay,
		parsedStartTime,
		parsedEndTime,
	)

	created, err := u.schRepo.CreateSchedule(ctx, schedule)
	if err != nil {
		if errors.Is(err, domain.ErrScheduleAlreadyExists) {
			return nil, err
		}

		u.log.Error("failed to create schedule", "layer", ScheduleLayer, "error", err)
		return nil, domain.ErrInternal
	}

	err = u.slots.GenerateSlotForSchedule(ctx, created, created.Day.Truncate(time.Hour*24))
	if err != nil {
		u.log.Error("failed to generate slots for schedule", "layer", ScheduleLayer,
			"id", created.ID, "error", err)
	}

	return output.NewCreateScheduleOutput(
		created.ID,
		created.RoomID,
		created.Day.Format(domain.DateLayout),
		minutesToTime(created.StartWorkTime),
		minutesToTime(created.EndWorkTime),
		created.CreatedAt.UTC(),
	), nil
}

func minutesToTime(minutes int) string {
	h := minutes / 60
	m := minutes % 60

	return fmt.Sprintf("%02d:%02d", h, m)
}
