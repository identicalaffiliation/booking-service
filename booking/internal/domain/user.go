package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Nickname     string    `json:"nickname"`
	PasswordHash string    `json:"-"`
	Role         UserRole  `json:"role"`
	CreatedAt    time.Time `json:"createdAt"`
}

func NewUser(nickname, passwordHash, role string) *User {
	return &User{
		ID:           uuid.New(),
		Nickname:     nickname,
		PasswordHash: passwordHash,
		Role:         UserRole(role),
	}
}
