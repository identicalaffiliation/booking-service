package ports

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/identicalaffiliation/booking-service/booking/internal/domain"
)

type Logger interface {
	Debug(msg string, args ...any)
	Error(msg string, args ...any)
}

type RoomsRepository interface {
	CreateRoom(ctx context.Context, r *domain.Room) (*domain.Room, error)
	GetRoom(ctx context.Context, id uuid.UUID) (*domain.Room, error)
	GetRooms(ctx context.Context) ([]*domain.Room, error)
	DeleteRoom(ctx context.Context, id uuid.UUID) error
}

type SchedulesRepository interface {
	CreateSchedule(ctx context.Context, s *domain.Schedule) (*domain.Schedule, error)
	GetAllSchedulesByToday(ctx context.Context, date time.Time) ([]*domain.Schedule, error)
}

type SlotsRepository interface {
	CreateSlot(ctx context.Context, slot *domain.Slot) error
}

type UsersRepository interface {
	CreateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	GetUser(ctx context.Context, nickname string) (*domain.User, error)
}

type RefreshTokensRepository interface {
	CreateRefreshToken(ctx context.Context, t *domain.RefreshToken) (*domain.RefreshToken, error)
	GetForUpdate(ctx context.Context, id uuid.UUID) (*domain.RefreshToken, error)
	Revoked(ctx context.Context, id uuid.UUID) error
}

type Hasher interface {
	Hash(val string) string
	CompareHash(hash, val string) bool
	HashPassword(password string) (string, error)
	ComparePassword(password, reqPassword string) error
}
