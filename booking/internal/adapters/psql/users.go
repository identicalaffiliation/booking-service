package psql

import (
	"context"
	"errors"
	"fmt"

	"github.com/identicalaffiliation/booking-service/booking/internal/domain"
	"github.com/jackc/pgx/v5"
)

type UsersRepository struct {
	db DBTX
}

func NewUsersRepository(db DBTX) *UsersRepository {
	return &UsersRepository{db: db}
}

func (r *UsersRepository) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	const query = `INSERT INTO users (id, nickname, password_hash, role) VALUES ($1, $2, $3, $4)
		RETURNING id, nickname, password_hash, role, created_at`

	var created domain.User
	err := r.db.QueryRow(
		ctx,
		query,
		user.ID,
		user.Nickname,
		user.PasswordHash,
		user.Role,
	).Scan(
		&created.ID,
		&created.Nickname,
		&created.PasswordHash,
		&created.Role,
		&created.CreatedAt,
	)
	if err != nil {
		if checkUniqueViolation(err) {
			return nil, domain.ErrUserAlreadyExists
		}

		return nil, fmt.Errorf("create user: %w", err)
	}

	return &created, nil
}

func (r *UsersRepository) GetUser(ctx context.Context, nickname string) (*domain.User, error) {
	const query = `SELECT id, nickname, password_hash, role, created_at FROM users
		WHERE nickname = $1`

	var user domain.User
	err := r.db.QueryRow(
		ctx,
		query,
		nickname,
	).Scan(
		&user.ID,
		&user.Nickname,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}

		return nil, fmt.Errorf("find user: %w", err)
	}

	return &user, nil
}
