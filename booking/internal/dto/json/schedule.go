package json

import (
	"time"

	"github.com/google/uuid"
)

type CreateScheduleRequest struct {
	Day   string `json:"day"`
	Start string `json:"start"`
	End   string `json:"end"`
}

type CreateScheduleResponse struct {
	ID        uuid.UUID `json:"id"`
	RoomID    uuid.UUID `json:"roomId"`
	Day       string    `json:"day"`
	Start     int       `json:"start"`
	End       int       `json:"end"`
	CreatedAt time.Time `json:"createdAt"`
}
