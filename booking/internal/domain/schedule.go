package domain

import (
	"time"

	"github.com/google/uuid"
)

type Schedule struct {
	ID            uuid.UUID
	RoomID        uuid.UUID
	Day           time.Time
	StartWorkTime time.Time
	EndWorkTime   time.Time
	CreatedAt     time.Time
}
