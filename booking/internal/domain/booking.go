package domain

import (
	"time"

	"github.com/google/uuid"
)

type Booking struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	SlotID    uuid.UUID
	Status    BookingStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}
