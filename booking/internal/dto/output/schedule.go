package output

import (
	"time"

	"github.com/google/uuid"
)

type CreateScheduleOutput struct {
	ID        uuid.UUID
	RoomID    uuid.UUID
	Day       string
	Start     string
	End       string
	CreatedAt time.Time
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
