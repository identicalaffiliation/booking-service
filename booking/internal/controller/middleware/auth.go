package middlewares

import (
	"context"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/identicalaffiliation/booking-service/booking/internal/domain"
	"github.com/identicalaffiliation/booking-service/booking/internal/dto/output"
	"github.com/labstack/echo/v4"
)

const (
	AuthorizationHeaderKey = "Authorization"
	BearerMethod           = "Bearer "
	userIdKey              = "userId"
	userRoleKey            = "role"
)

func AuthMiddleware(secret, issuer string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			header := ctx.Request().Header
			method := header.Get(AuthorizationHeaderKey)
			if method == "" {
				return output.NewNotAuthorized("invalid header")
			}

			if !strings.HasPrefix(method, BearerMethod) {
				return output.NewNotAuthorized("invalid auth method")
			}

			tokenString := strings.TrimPrefix(method, BearerMethod)
			if tokenString == "" {
				return output.NewNotAuthorized("invalid token type")
			}

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("get signed method: %v", token.Header["alg"])
				}

				return []byte(secret), nil
			})
			if err != nil || !token.Valid {
				return output.NewNotAuthorized("invalid token")
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return output.NewNotAuthorized("invalid token claims")
			}

			iss, ok := claims["iss"].(string)
			if !ok || iss != issuer {
				return output.NewNotAuthorized("invalid token claims")
			}

			id, ok := claims["sub"].(string)
			if !ok || id == "" {
				return output.NewNotAuthorized("invalid token claims")
			}

			roleStr, ok := claims["role"].(string)
			if !ok || roleStr == "" {
				return output.NewNotAuthorized("invalid token claims")
			}

			role := domain.UserRole(roleStr)

			if role != domain.Admin && role != domain.Client {
				return output.NewNotAuthorized("invalid token claims")
			}

			reqCtx := ctx.Request().Context()
			reqCtx = context.WithValue(reqCtx, userIdKey, id)
			reqCtx = context.WithValue(reqCtx, userRoleKey, role)

			ctx.SetRequest(ctx.Request().WithContext(reqCtx))

			return next(ctx)
		}
	}
}
