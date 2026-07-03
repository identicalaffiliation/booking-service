package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/identicalaffiliation/booking-service/booking/config"
	"github.com/identicalaffiliation/booking-service/booking/internal/domain"
	"github.com/identicalaffiliation/booking-service/booking/internal/ports"
)

const SlotsLayer = "slots service layer"

type SlotsUsecase struct {
	slotsRepo     ports.SlotsRepository
	schedulesRepo ports.SchedulesRepository
	log           ports.Logger
	cfg           *config.BookingConfig
}

func NewSlotsUsecase(
	slotsRepo ports.SlotsRepository,
	schedulesRepo ports.SchedulesRepository,
	log ports.Logger,
	cfg *config.BookingConfig,
) *SlotsUsecase {
	return &SlotsUsecase{
		slotsRepo:     slotsRepo,
		schedulesRepo: schedulesRepo,
		log:           log,
		cfg:           cfg,
	}
}

func (u *SlotsUsecase) GenerateSlots(ctx context.Context) error {
	date := time.Now().UTC().Truncate(time.Hour * 24)
	availableSchedules, err := u.schedulesRepo.GetAllSchedulesByToday(ctx, date)
	if err != nil {
		return domain.ErrInternal
	}

	for _, schedule := range availableSchedules {
		err := u.GenerateSlotForSchedule(ctx, schedule, date)
		if err != nil {
			u.log.Error("failed to generate slot", "layer", SlotsLayer, "error", err)
			continue
		}
	}

	return nil
}

func (u *SlotsUsecase) GenerateSlotForSchedule(ctx context.Context, schedule *domain.Schedule, date time.Time) error {
	interval := int(u.cfg.SlotInterval.Minutes())
	for i := schedule.StartWorkTime; i+interval <= schedule.EndWorkTime; i += interval {
		slot := domain.NewSlot(schedule.RoomID, date, i, i+interval)
		err := u.slotsRepo.CreateSlot(ctx, slot)
		if err != nil {
			if errors.Is(err, domain.ErrSlotAlreadyExists) {
				continue
			}

			return err
		}
	}

	return nil
}
