package input

import (
	"github.com/google/uuid"
)

type CreateScheduleInput struct {
	RoomID uuid.UUID `json:"roomId"`
	Day    string    `json:"day"`
	Start  string    `json:"start"`
	End    string    `json:"end"`
}

func NewCreateScheduleInput(roomID uuid.UUID, day, start, end string) *CreateScheduleInput {
	return &CreateScheduleInput{
		RoomID: roomID,
		Day:    day,
		Start:  start,
		End:    end,
	}
}
