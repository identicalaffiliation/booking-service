package generator

import (
	"context"
	"time"

	"github.com/identicalaffiliation/booking-service/booking/config"
	"github.com/identicalaffiliation/booking-service/booking/internal/ports"
	"github.com/identicalaffiliation/booking-service/booking/internal/usecase"
)

const (
	seconds     = 0
	nanoseconds = 0
)

func StartSlotGenerator(
	ctx context.Context,
	service *usecase.SlotsUsecase,
	log ports.Logger,
	cfg *config.BookingConfig) {
	go func() {
		for {
			now := time.Now().UTC()
			nextRound := getNextRound(now, cfg)

			timer := time.NewTimer(time.Until(nextRound))
			select {
			case <-ctx.Done():
				log.Debug("slot generator is stopped by ctx")
				timer.Stop()
				return
			case <-timer.C:
				log.Debug("slot generator started..")
				if err := service.GenerateSlots(ctx); err != nil {
					log.Error("slot generation fail", "error", err)
				}
			}
		}
	}()
}

func getNextRound(now time.Time, cfg *config.BookingConfig) time.Time {
	next := time.Date(now.Year(), now.Month(), now.Day(),
		cfg.JobStartHours, cfg.JobStartMinutes, seconds, nanoseconds,
		time.UTC,
	)

	if now.After(next) {
		next = next.Add(time.Hour * 24)
	}

	return next
}
