package output

import "github.com/identicalaffiliation/booking-service/booking/internal/domain"

type CreateBookingOutput struct {
	Booking *domain.Booking `json:"booking"`
}

func NewBookingOutput(b *domain.Booking) *CreateBookingOutput {
	return &CreateBookingOutput{Booking: b}
}
