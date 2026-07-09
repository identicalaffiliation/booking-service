package domain

import (
	"time"

	"github.com/google/uuid"
)

type Booking struct {
	ID        uuid.UUID     `json:"id"`
	UserID    uuid.UUID     `json:"userId"`
	SlotID    uuid.UUID     `json:"slotId"`
	Status    BookingStatus `json:"status"`
	CreatedAt time.Time     `json:"createdAt"`
	UpdatedAt time.Time     `json:"updatedAt"`
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
