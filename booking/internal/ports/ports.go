package ports

import (
	"context"

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
}
