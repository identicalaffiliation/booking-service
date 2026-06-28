package application

import (
	"context"
	"errors"
	"time"

	"github.com/identicalaffiliation/booking-service/booking/internal/adapters/storage/psql"
	"github.com/identicalaffiliation/booking-service/booking/internal/config"
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
	availableSchedules, err := u.schedulesRepo.GetAllSchedules(ctx)
	if err != nil {
		return ErrInternal
	}

	for _, schedule := range availableSchedules {
		err := u.GenerateSlotForSchedule(ctx, schedule)
		if err != nil {
			u.log.Error("failed to generate slot", "layer", SlotsLayer, "error", err)
			continue
		}
	}

	return nil
}

func (u *SlotsUsecase) GenerateSlotForSchedule(ctx context.Context, schedule *domain.Schedule) error {
	interval := int(u.cfg.SlotInterval.Minutes())
	date := schedule.Day.Truncate(time.Hour * 24)
	for i := schedule.StartWorkTime; i+interval <= schedule.EndWorkTime; i += interval {
		slot := domain.NewSlot(schedule.RoomID, date, i, i+interval)
		err := u.slotsRepo.CreateSlot(ctx, slot)
		if err != nil {
			if errors.Is(err, psql.ErrSlotAlreadyExists) {
				continue
			}

			return err
		}
	}

	return nil
}
