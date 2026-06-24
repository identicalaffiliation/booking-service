package ports

import (
	"context"

	"github.com/identicalaffiliation/booking-service/booking/internal/domain"
)

type RoomsRepository interface {
	CreateRoom(ctx context.Context, r *domain.Room) (*domain.Room, error)
}
