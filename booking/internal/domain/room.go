package domain

import (
	"time"

	"github.com/google/uuid"
)

type Room struct {
	ID        uuid.UUID
	Name      string
	Capacity  int
	CreatedAt time.Time
}

func NewRoom(name string, capacity int) *Room {
	return &Room{
		ID:       uuid.New(),
		Name:     name,
		Capacity: capacity,
	}
}
