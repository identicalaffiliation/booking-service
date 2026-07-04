package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/identicalaffiliation/booking-service/booking/config"
	"github.com/identicalaffiliation/booking-service/booking/internal/adapters/psql"
	"github.com/identicalaffiliation/booking-service/booking/internal/domain"
	"github.com/identicalaffiliation/booking-service/booking/internal/dto/input"
	"github.com/identicalaffiliation/booking-service/booking/internal/dto/output"
	"github.com/identicalaffiliation/booking-service/booking/internal/ports"
	"github.com/identicalaffiliation/booking-service/booking/pkg/hasher"
)

type AuthUsecase struct {
	hasher     ports.Hasher
	log        ports.Logger
	usersRepo  ports.UsersRepository
	tokensRepo ports.RefreshTokensRepository
	cfg        *config.BookingConfig
	txManager  *psql.TxManager
}

func NewAuthUsecase(
	users ports.UsersRepository,
	tokens ports.RefreshTokensRepository,
	log ports.Logger,
	cfg *config.BookingConfig,
	txManager *psql.TxManager,
) *AuthUsecase {
	return &AuthUsecase{
		hasher:     hasher.NewHasher(),
		usersRepo:  users,
		tokensRepo: tokens,
		log:        log,
		cfg:        cfg,
		txManager:  txManager,
	}
}

func (u *AuthUsecase) Registration(ctx context.Context, in *input.CreateUserInput) (*output.UserOutput, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	hashedPassword, err := u.hasher.HashPassword(in.Password)
	if err != nil {
		u.log.Error("failed to hash password", "error", err)
		return nil, domain.ErrInternal
	}

	user := domain.NewUser(in.Nickname, hashedPassword, in.Role)
	created, err := u.usersRepo.CreateUser(ctx, user)
	if err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			return nil, err
		}

		u.log.Error("failed to create user", "error", err)
		return nil, domain.ErrInternal
	}

	return output.NewUserOutput(created.ID, created.Nickname, created.Role, created.CreatedAt), nil
}

func (u *AuthUsecase) Login(ctx context.Context, in *input.LoginInput) (*output.LoginOutput, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	user, err := u.usersRepo.GetUser(ctx, in.Nickname)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, domain.ErrInvalidUserData
		}

		u.log.Error("failed to get user", "nickname", in.Nickname, "error", err)
		return nil, domain.ErrInternal
	}

	if err := u.hasher.ComparePassword(user.PasswordHash, in.Password); err != nil {
		return nil, domain.ErrInvalidUserData
	}

	accessToken, err := u.generateAccessToken(user.ID, string(user.Role))
	if err != nil {
		return nil, err
	}

	now := time.Now()
	tokenID := uuid.New()

	rawToken, err := u.generateRefreshToken(tokenID, user.ID, now)
	if err != nil {
		return nil, err
	}

	token := domain.NewRefreshToken(
		tokenID,
		user.ID,
		u.hasher.Hash(rawToken),
		false,
		now.Add(u.cfg.RefreshTokenConfig.ExpiredAt).Unix(),
	)

	_, err = u.tokensRepo.CreateRefreshToken(ctx, token)
	if err != nil {
		return nil, domain.ErrInternal
	}

	return output.NewLoginOutput(accessToken, rawToken), nil
}

func (u *AuthUsecase) Refresh(ctx context.Context, in *input.RefreshAccessTokenInput) (*output.Tokens, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	var out output.Tokens
	err := u.txManager.WithTx(ctx, func(ctx context.Context, tx psql.DBTX) error {
		claims, err := u.parseRefreshToken(in.RefreshToken)
		if err != nil {
			return domain.ErrInvalidRefreshTokenData
		}

		if err := u.validateRefreshToken(claims); err != nil {
			return domain.ErrInvalidRefreshTokenData
		}

		tokenID, err := uuid.Parse(claims["sub"].(string))
		if err != nil {
			return domain.ErrInvalidRefreshTokenData
		}

		userID, err := uuid.Parse(claims["userId"].(string))
		if err != nil {
			return domain.ErrInvalidRefreshTokenData
		}

		oldRefreshToken, err := u.tokensRepo.GetForUpdate(ctx, tokenID)
		if err != nil {
			if errors.Is(err, domain.ErrTokenNotFound) {
				return err
			}

			u.log.Error("failed to get refresh old token", "error", err)
			return domain.ErrInternal
		}

		if !u.hasher.CompareHash(in.RefreshToken, oldRefreshToken.TokenHash) {
			return domain.ErrInvalidRefreshTokenData
		}

		if err := u.tokensRepo.Revoked(ctx, oldRefreshToken.ID); err != nil {
			if errors.Is(err, domain.ErrTokenNotFound) {
				return err
			}

			u.log.Error("failed to update refresh token", "error", err)
			return domain.ErrInternal
		}

		now := time.Now()
		tokenID = uuid.New()

		rawToken, err := u.generateRefreshToken(tokenID, userID, now)
		if err != nil {
			return err
		}

		newHashedToken := domain.NewRefreshToken(
			tokenID,
			userID,
			u.hasher.Hash(rawToken),
			true,
			now.Add(u.cfg.RefreshTokenConfig.ExpiredAt).Unix(),
		)

		_, err = u.tokensRepo.CreateRefreshToken(ctx, newHashedToken)
		if err != nil {
			return domain.ErrInternal
		}

		out.RefreshToken = rawToken
		out.AccessToken = ""

		return nil
	})
	if err != nil {
		return nil, err
	}
	
	return &out, nil
}

func (u *AuthUsecase) generateAccessToken(id uuid.UUID, role string) (string, error) {
	now := time.Now()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  id.String(),
		"role": role,
		"iss":  u.cfg.AccessTokenConfig.IssuedBy,
		"exp":  now.Add(u.cfg.AccessTokenConfig.ExpiredAt).Unix(),
		"iat":  now.Unix(),
	})

	raw, err := token.SignedString([]byte(u.cfg.JwtSecret))
	if err != nil {
		return "", domain.ErrInternal
	}

	return raw, nil
}

func (u *AuthUsecase) generateRefreshToken(tokenId, userId uuid.UUID, now time.Time) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":    tokenId.String(),
		"userId": userId.String(),
		"iss":    u.cfg.RefreshTokenConfig.IssuedBy,
		"exp":    now.Add(u.cfg.RefreshTokenConfig.ExpiredAt).Unix(),
		"iat":    now.Unix(),
	})

	raw, err := token.SignedString([]byte(u.cfg.JwtSecret))
	if err != nil {
		return "", domain.ErrInternal
	}

	return raw, nil
}

func (u *AuthUsecase) parseRefreshToken(token string) (jwt.MapClaims, error) {
	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("get signed method: %v\n", token.Header["alg"])
		}

		return []byte(u.cfg.JwtSecret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("parse jwt token: %w", err)
	}

	return jwtToken.Claims.(jwt.MapClaims), nil
}

func (u *AuthUsecase) validateRefreshToken(claims jwt.MapClaims) error {
	sub, ok := claims["sub"].(string)
	if !ok || sub == "" {
		return domain.ErrInvalidRefreshTokenData
	}

	userId, ok := claims["userId"].(string)
	if !ok || userId == "" {
		return domain.ErrInvalidRefreshTokenData
	}

	iss, ok := claims["iss"].(string)
	if !ok || iss != u.cfg.RefreshTokenConfig.IssuedBy {
		return domain.ErrInvalidRefreshTokenData
	}

	return nil
}
