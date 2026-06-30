package input

import (
	"github.com/google/uuid"
	"github.com/identicalaffiliation/booking-service/booking/internal/domain"
	"github.com/identicalaffiliation/booking-service/booking/pkg/validator"
)

type CreateScheduleInput struct {
	RoomID uuid.UUID `json:"roomId" validate:"required"`
	Day    string    `json:"day" validate:"required,datetime=2006-01-02"`
	Start  string    `json:"start" validate:"required,datetime=15:04"`
	End    string    `json:"end" validate:"required,datetime=15:04"`
}

func (in *CreateScheduleInput) Validate() error {
	if err := validator.NewValidator().Validate(in); err != nil {
		return domain.ErrInvalidScheduleData
	}

	if in.RoomID == uuid.Nil {
		return domain.ErrInvalidScheduleData
	}

	_, err := domain.ParseTimeDate(in.Day)
	if err != nil {
		return domain.ErrInvalidScheduleData
	}

	start, err := domain.ParseTimeDuration(in.Start)
	if err != nil {
		return domain.ErrInvalidScheduleData
	}

	end, err := domain.ParseTimeDuration(in.End)
	if err != nil {
		return domain.ErrInvalidScheduleData
	}

	if start >= end {
		return domain.ErrInvalidScheduleData
	}
	
	if (start%60 != 0) || (end%60 != 0) {
		return domain.ErrInvalidScheduleData
	}

	return nil
}
