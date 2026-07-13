package ports

import (
	"context"

	"github.com/google/uuid"
)

type Logger interface {
	Debug(msg string, args ...any)
	Error(msg string, args ...any)
}

type Repository interface {
	Create(ctx context.Context, id uuid.UUID, payload []byte) error
}
