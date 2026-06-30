package output

import (
	"time"

	"github.com/google/uuid"
)

type CreateScheduleOutput struct {
	ID        uuid.UUID `json:"id"`
	RoomID    uuid.UUID `json:"roomId"`
	Day       string    `json:"day"`
	Start     string    `json:"start"`
	End       string    `json:"end"`
	CreatedAt time.Time `json:"createdAt"`
}

func NewCreateScheduleOutput(
	id, roomId uuid.UUID,
	day string,
	start, end string,
	createdAt time.Time,
) *CreateScheduleOutput {
	return &CreateScheduleOutput{
		ID:        id,
		RoomID:    roomId,
		Day:       day,
		Start:     start,
		End:       end,
		CreatedAt: createdAt.UTC(),
	}
}
