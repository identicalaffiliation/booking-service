package ports

import (
	"context"
	"time"

	"github.com/identicalaffiliation/booking-service/booking/internal/domain"
)

type Logger interface {
	Debug(msg string, args ...any)
	Error(msg string, args ...any)
}

type RoomsRepository interface {
	CreateRoom(ctx context.Context, r *domain.Room) (*domain.Room, error)
}

type SchedulesRepository interface {
	CreateSchedule(ctx context.Context, s *domain.Schedule) (*domain.Schedule, error)
	GetAllSchedules(ctx context.Context, date time.Time) ([]*domain.Schedule, error)
}

type SlotsRepository interface {
	CreateSlot(ctx context.Context, slot *domain.Slot) error
}
