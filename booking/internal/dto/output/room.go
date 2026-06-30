package output

import (
	"time"

	"github.com/google/uuid"
)

type CreateRoomOutput struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Capacity  int       `json:"capacity"`
	CreatedAt time.Time `json:"createdAt"`
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
