package domain

import (
	"time"

	"github.com/google/uuid"
)

type Room struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Capacity  int       `json:"capacity"`
	CreatedAt time.Time `json:"createdAt"`
}

func NewRoom(name string, capacity int) *Room {
	return &Room{
		ID:       uuid.New(),
		Name:     name,
		Capacity: capacity,
	}
}
