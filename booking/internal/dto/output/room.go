package output

import (
	"time"

	"github.com/google/uuid"
)

type CreateRoomOutput struct {
	ID        uuid.UUID
	Name      string
	Capacity  int
	CreatedAt time.Time
}

func NewCreateRoomOutput(
	id uuid.UUID,
	name string,
	capacity int,
	createdTime time.Time,
) *CreateRoomOutput {
	return &CreateRoomOutput{
		ID:        id,
		Name:      name,
		Capacity:  capacity,
		CreatedAt: createdTime.UTC(),
	}
}
