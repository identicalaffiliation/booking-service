package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	Nickname     string
	PasswordHash string
	Role         UserRole
	CreatedAt    time.Time
}
