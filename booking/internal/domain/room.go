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
