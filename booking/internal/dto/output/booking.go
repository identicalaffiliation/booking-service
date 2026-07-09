package output

import (
	"time"

	"github.com/google/uuid"
	"github.com/identicalaffiliation/booking-service/booking/internal/domain"
)

type CreateBookingOutput struct {
	Booking *domain.Booking `json:"booking"`
}

func NewBookingOutput(b *domain.Booking) *CreateBookingOutput {
	return &CreateBookingOutput{Booking: b}
}

type MyBookingsOutput struct {
	Bookings []*Booking `json:"bookings"`
}

type MyBookingOutput struct {
	Booking *MyBooking `json:"booking"`
}

type MyBooking struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"userId"`

	SlotID        uuid.UUID `json:"slotId"`
	SlotDate      string    `json:"slotDate"`
	SlotStartTime string    `json:"slotStartTime"`
	SlotEndTime   string    `json:"slotEndTime"`

	RoomID       uuid.UUID `json:"roomId"`
	RoomName     string    `json:"roomName"`
	RoomCapacity int       `json:"roomCapacity"`

	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Booking struct {
	Booking *domain.Booking `json:"booking"`
}

func NewMyBookingOutput(
	bookingID, userID, slotID, roomID uuid.UUID,
	slotDate, slotStartTime, slotEndTime string,
	roomName string, roomCapacity int, bookingStatus string,
	createdAt, updatedAt time.Time,
) *MyBookingOutput {
	booking := &MyBooking{
		ID:            bookingID,
		UserID:        userID,
		SlotID:        slotID,
		RoomID:        roomID,
		SlotDate:      slotDate,
		SlotStartTime: slotStartTime,
		SlotEndTime:   slotEndTime,
		RoomName:      roomName,
		RoomCapacity:  roomCapacity,
		Status:        bookingStatus,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}

	return &MyBookingOutput{
		Booking: booking,
	}
}

func NewMyBookingsOutput(bookings []*domain.Booking) *MyBookingsOutput {
	var output MyBookingsOutput
	for _, booking := range bookings {
		out := &Booking{
			Booking: booking,
		}
		output.Bookings = append(output.Bookings, out)
	}

	return &output
}
