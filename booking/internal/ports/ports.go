package ports

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/identicalaffiliation/booking-service/booking/internal/domain"
	"github.com/identicalaffiliation/booking-service/booking/internal/dto/input"
	"github.com/identicalaffiliation/booking-service/booking/internal/dto/output"
)

type Logger interface {
	Debug(msg string, args ...any)
	Error(msg string, args ...any)
}

type BookingsRepository interface {
	CreateBooking(ctx context.Context, booking *domain.Booking) (*domain.Booking, error)
	GetMyBooking(ctx context.Context, bookingID uuid.UUID) (*domain.MyBooking, error)
	GetMyBookings(ctx context.Context, userID uuid.UUID) ([]*domain.Booking, error)
	CancelMyBooking(ctx context.Context, bookingID uuid.UUID) error
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
	GetForUpdateRefreshToken(ctx context.Context, id uuid.UUID) (*domain.RefreshToken, error)
	Revoked(ctx context.Context, id uuid.UUID) error
}

type Hasher interface {
	Hash(val string) string
	CompareHash(hash, val string) bool
	HashPassword(password string) (string, error)
	ComparePassword(password, reqPassword string) error
}

type AuthUsecase interface {
	Registration(ctx context.Context, in *input.CreateUserInput) (*output.UserOutput, error)
	Login(ctx context.Context, in *input.LoginInput) (*output.LoginOutput, error)
	Refresh(ctx context.Context, in *input.RefreshAccessTokenInput) (*output.LoginOutput, error)
}

type RoomsUsecase interface {
	CreateRoom(ctx context.Context, in *input.CreateRoomInput) (*output.CreateRoomOutput, error)
	GetRoom(ctx context.Context, id string) (*output.RoomOutput, error)
	GetRooms(ctx context.Context) (*output.RoomsOutput, error)
	DeleteRoom(ctx context.Context, id string) error
}

type SchedulesUsecase interface {
	CreateSchedule(ctx context.Context, in *input.CreateScheduleInput) (*output.CreateScheduleOutput, error)
}

type BookingsUsecase interface {
	Create(ctx context.Context, in *input.CreateBookingInput) (*output.CreateBookingOutput, error)
	Cancel(ctx context.Context, bookingID uuid.UUID) error
}
