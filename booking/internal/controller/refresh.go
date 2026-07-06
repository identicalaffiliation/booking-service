package controller

import (
	"errors"
	"net/http"

	"github.com/identicalaffiliation/booking-service/booking/internal/domain"
	"github.com/identicalaffiliation/booking-service/booking/internal/dto/input"
	"github.com/identicalaffiliation/booking-service/booking/internal/dto/output"
	"github.com/identicalaffiliation/booking-service/booking/internal/usecase"
	"github.com/labstack/echo/v4"
)

func Refresh(auth *usecase.AuthUsecase) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var in input.RefreshAccessTokenInput
		if err := ctx.Bind(&in); err != nil {
			return output.NewBadRequest("invalid json body")
		}

		reqCtx := ctx.Request().Context()

		out, err := auth.Refresh(reqCtx, &in)
		if err != nil {
			switch {
			case errors.Is(err, domain.ErrInvalidRefreshTokenData):
				return output.NewBadRequest("invalid token data")
			case errors.Is(err, domain.ErrTokenNotFound):
				return output.NewNotFound("token not found")
			default:
				return output.NewInternal()
			}
		}

		return ctx.JSON(http.StatusOK, out)
	}
}
