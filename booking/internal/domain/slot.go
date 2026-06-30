package domain

import (
	"time"

	"github.com/google/uuid"
)

type Slot struct {
	ID        uuid.UUID
	RoomID    uuid.UUID
	Day       time.Time
	StartTime int
	EndTime   int
	CreatedAt time.Time
}

func NewSlot(roomID uuid.UUID, day time.Time, startTime, endTime int) *Slot {
	return &Slot{
		ID:        uuid.New(),
		RoomID:    roomID,
		Day:       day,
		StartTime: startTime,
		EndTime:   endTime,
	}
}
