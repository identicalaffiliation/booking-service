package output

import (
	"time"

	"github.com/google/uuid"
	"github.com/identicalaffiliation/booking-service/booking/internal/domain"
)

type UserOutput struct {
	ID        uuid.UUID       `json:"id"`
	Nickname  string          `json:"nickname"`
	Role      domain.UserRole `json:"role"`
	CreatedAt time.Time       `json:"createdAt"`
}

func NewUserOutput(id uuid.UUID, nickname string, role domain.UserRole, created time.Time) *UserOutput {
	return &UserOutput{
		ID:        id,
		Nickname:  nickname,
		Role:      role,
		CreatedAt: created,
	}
}
