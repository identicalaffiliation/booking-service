package domain

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	TokenHash string
	ExpiresAt int64
	Revoked   bool
	CreatedAt time.Time
}

func NewRefreshToken(
	tokenID uuid.UUID,
	userID uuid.UUID,
	tokenHash string,
	revoked bool,
	expAt int64,
) *RefreshToken {
	return &RefreshToken{
		ID:        tokenID,
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: expAt,
		Revoked:   revoked,
	}
}
