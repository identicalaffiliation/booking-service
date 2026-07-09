package domain

import "github.com/google/uuid"

type KafkaEvent struct {
	BookingID string `json:"bookingId"`
	UserID    string `json:"userId"`
	SlotID    string `json:"slotId"`
	EventType string `json:"type"`
}

func NewEvent(bookingID, userID, slotID uuid.UUID, eventType string) *KafkaEvent {
	return &KafkaEvent{
		BookingID: bookingID.String(),
		UserID:    userID.String(),
		SlotID:    slotID.String(),
		EventType: eventType,
	}
}
