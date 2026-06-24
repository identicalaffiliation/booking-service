package json

import (
	"time"

	"github.com/google/uuid"
)

type CreateRoomRequest struct {
	Name     string `json:"name"`
	Capacity int    `json:"capacity"`
}

type CreateRoomResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Capacity  int       `json:"capacity"`
	CreatedAt time.Time `json:"createdAt"`
}
