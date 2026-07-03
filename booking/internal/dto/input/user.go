package input

import (
	"strings"

	"github.com/identicalaffiliation/booking-service/booking/internal/domain"
	"github.com/identicalaffiliation/booking-service/booking/pkg/validator"
)

type CreateUserInput struct {
	Nickname string `json:"nickname" validate:"required,min=1,max=50"`
	Password string `json:"password" validate:"required,min=8"`
	Role     string `json:"role" validate:"required,oneof=client admin"`
}

func (in *CreateUserInput) Validate() error {
	if err := validator.NewValidator().Validate(in); err != nil {
		return domain.ErrInvalidUserData
	}

	if len(strings.TrimSpace(in.Nickname)) < 1 {
		return domain.ErrInvalidUserData
	}

	if len(strings.TrimSpace(in.Password)) < 8 {
		return domain.ErrInvalidUserData
	}

	return nil
}

type LoginInput struct {
	Nickname string `json:"nickname" validate:"required,min=1,max=50"`
	Password string `json:"password" validate:"required,min=8"`
}

func (in *LoginInput) Validate() error {
	if err := validator.NewValidator().Validate(in); err != nil {
		return domain.ErrInvalidUserData
	}

	if len(strings.TrimSpace(in.Nickname)) < 1 {
		return domain.ErrInvalidUserData
	}

	if len(strings.TrimSpace(in.Password)) < 8 {
		return domain.ErrInvalidUserData
	}

	return nil
}
