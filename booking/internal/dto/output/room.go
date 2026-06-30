package output

import (
	"time"

	"github.com/google/uuid"
	"github.com/identicalaffiliation/booking-service/booking/internal/domain"
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

type RoomOutput struct {
	Room domain.Room `json:"room"`
}

func NewRoomOutput(id uuid.UUID, name string, cap int, created time.Time) *RoomOutput {
	return &RoomOutput{
		Room: domain.Room{
			ID:        id,
			Name:      name,
			Capacity:  cap,
			CreatedAt: created,
		},
	}
}

type RoomsOutput struct {
	Rooms []*RoomOutput `json:"rooms"`
}
