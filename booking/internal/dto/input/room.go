package input

import (
	"strings"

	"github.com/identicalaffiliation/booking-service/booking/internal/domain"
	"github.com/identicalaffiliation/booking-service/booking/pkg/validator"
)

type CreateRoomInput struct {
	Name     string `json:"name" validate:"required,min=6,max=50"`
	Capacity int    `json:"capacity" validate:"required,min=1,max=7"`
}

func (in *CreateRoomInput) Validate() error {
	if err := validator.NewValidator().Validate(in); err != nil {
		return domain.ErrInvalidRoomData
	}

	if len(strings.TrimSpace(in.Name)) == 0 {
		return domain.ErrInvalidRoomData
	}

	return nil
}
