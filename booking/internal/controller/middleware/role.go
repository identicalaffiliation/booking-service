package middlewares

import (
	"github.com/identicalaffiliation/booking-service/booking/internal/domain"
	"github.com/identicalaffiliation/booking-service/booking/internal/dto/output"
	"github.com/labstack/echo/v4"
)

func RoleMiddleware(requiredRole domain.UserRole) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			if requiredRole != ctx.Request().Context().Value(userRoleKey).(domain.UserRole) {
				return output.NewForbidden("access denied")
			}

			return next(ctx)
		}
	}
}
