package input

import (
	"github.com/google/uuid"
	"github.com/identicalaffiliation/booking-service/booking/internal/domain"
	"github.com/identicalaffiliation/booking-service/booking/pkg/validator"
)

type CreateBookingInput struct {
	SlotID uuid.UUID `json:"slotId" validate:"required"`
}

func (in *CreateBookingInput) Validate() error {
	if err := validator.NewValidator().Validate(in); err != nil {
		return domain.ErrInvalidBookingData
	}

	if in.SlotID == uuid.Nil {
		return domain.ErrInvalidBookingData
	}

	return nil
}
