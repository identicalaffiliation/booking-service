package domain

import (
	"time"

	"github.com/google/uuid"
)

type Schedule struct {
	ID            uuid.UUID
	RoomID        uuid.UUID
	Day           time.Time
	StartWorkTime int
	EndWorkTime   int
	CreatedAt     time.Time
}

func NewSchedule(roomID uuid.UUID, workDay time.Time, startWorkTime, endWorkTime int) *Schedule {
	return &Schedule{
		ID:            uuid.New(),
		RoomID:        roomID,
		Day:           workDay,
		StartWorkTime: startWorkTime,
		EndWorkTime:   endWorkTime,
	}
}
