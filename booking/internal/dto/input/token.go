package input

import (
	"strings"

	"github.com/identicalaffiliation/booking-service/booking/internal/domain"
	"github.com/identicalaffiliation/booking-service/booking/pkg/validator"
)

type RefreshAccessTokenInput struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

func (in *RefreshAccessTokenInput) Validate() error {
	if err := validator.NewValidator().Validate(in); err != nil {
		return domain.ErrInvalidRefreshTokenData
	}

	if len(strings.TrimSpace(in.RefreshToken)) == 0 {
		return domain.ErrInvalidRefreshTokenData
	}

	return nil
}
