package input

import (
	"github.com/google/uuid"
)

type CreateScheduleInput struct {
	RoomID uuid.UUID
	Day    string
	Start  string
	End    string
}

func NewCreateScheduleInput(roomID uuid.UUID, day, start, end string) *CreateScheduleInput {
	return &CreateScheduleInput{
		RoomID: roomID,
		Day:    day,
		Start:  start,
		End:    end,
	}
}
