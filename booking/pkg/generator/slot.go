package generator

import (
	"context"
	"time"

	"github.com/identicalaffiliation/booking-service/booking/config"
	"github.com/identicalaffiliation/booking-service/booking/internal/ports"
	"github.com/identicalaffiliation/booking-service/booking/internal/usecase"
)

func StartSlotGenerator(
	ctx context.Context,
	service *usecase.SlotsUsecase,
	log ports.Logger,
	cfg *config.BookingConfig) {
	go func() {
		ticker := time.NewTicker(cfg.TickInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Debug("slot generator is stopped by ctx")
				return
			case <-ticker.C:
				log.Debug("slot generator started..")
				if err := service.GenerateSlots(ctx); err != nil {
					log.Error("slot generation fail", "error", err)
				}
			}
		}
	}()
}
