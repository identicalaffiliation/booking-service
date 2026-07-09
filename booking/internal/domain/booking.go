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

func NewBooking(userID, slotID uuid.UUID) *Booking {
	return &Booking{
		ID:     uuid.New(),
		UserID: userID,
		SlotID: slotID,
		Status: Active,
	}
}

type MyBooking struct {
	Booking       Booking
	SlotByBooking Slot
	RoomByBooking Room
}
