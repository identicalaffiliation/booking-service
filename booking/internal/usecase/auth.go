package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/identicalaffiliation/booking-service/booking/config"
	"github.com/identicalaffiliation/booking-service/booking/internal/domain"
	"github.com/identicalaffiliation/booking-service/booking/internal/dto/input"
	"github.com/identicalaffiliation/booking-service/booking/internal/dto/output"
	"github.com/identicalaffiliation/booking-service/booking/internal/ports"
	"github.com/identicalaffiliation/booking-service/booking/pkg/hasher"
)

type AuthUsecase struct {
	hasher ports.Hasher
	log    ports.Logger
	repo   ports.UsersRepository
	cfg    *config.BookingConfig
}

func NewAuthUsecase(repo ports.UsersRepository, log ports.Logger, cfg *config.BookingConfig) *AuthUsecase {
	return &AuthUsecase{
		hasher: hasher.NewHasher(),
		repo:   repo,
		log:    log,
		cfg:    cfg,
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
	created, err := u.repo.CreateUser(ctx, user)
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

	user, err := u.repo.GetUser(ctx, in.Nickname)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, domain.ErrInvalidUserData
		}

		u.log.Error("failed to get user", "nickname", in.Nickname, "error", err)
		return nil, domain.ErrInternal
	}

	if err := u.hasher.Compare(user.PasswordHash, in.Password); err != nil {
		return nil, domain.ErrInvalidUserData
	}

	accessToken, err := u.generateAccessToken(user.ID, string(user.Role))
	if err != nil {
		return nil, err
	}

	return output.NewLoginOutput(accessToken, "token"), nil
}

func (u *AuthUsecase) generateAccessToken(id uuid.UUID, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  id.String(),
		"role": role,
		"iss":  u.cfg.IssuedBy,
		"exp":  time.Now().Add(u.cfg.ExpiredAt).Unix(),
		"iat":  time.Now().Unix(),
	})

	raw, err := token.SignedString([]byte(u.cfg.JwtSecret))
	if err != nil {
		return "", domain.ErrInternal
	}

	return raw, nil
}

func (u *AuthUsecase) generateRefreshToken() (string, error) {
	//TODO
	return "", nil
}
