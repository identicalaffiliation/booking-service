package psql

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/identicalaffiliation/booking-service/booking/internal/domain"
	"github.com/jackc/pgx/v5"
)

type RefreshTokensRepository struct {
	db DBTX
}

func NewRefreshTokensRepository(db DBTX) *RefreshTokensRepository {
	return &RefreshTokensRepository{
		db: db,
	}
}

func (r *RefreshTokensRepository) CreateRefreshToken(
	ctx context.Context,
	t *domain.RefreshToken,
) (*domain.RefreshToken, error) {
	const query = `
		INSERT INTO 
		refresh_tokens (id, user_id, token_hash, expires_at, revoked) 
		VALUES 
		($1, $2, $3, $4, $5) 
		RETURNING 
		id, user_id, token_hash, expires_at, revoked, created_at`

	var created domain.RefreshToken
	err := r.db.QueryRow(
		ctx,
		query,
		t.ID,
		t.UserID,
		t.TokenHash,
		t.ExpiresAt,
		t.Revoked,
	).Scan(
		&created.ID,
		&created.UserID,
		&created.TokenHash,
		&created.ExpiresAt,
		&created.Revoked,
		&created.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create refresh token: %w", err)
	}

	return &created, nil
}

func (r *RefreshTokensRepository) GetForUpdateRefreshToken(
	ctx context.Context,
	id uuid.UUID,
) (*domain.RefreshToken, error) {
	const query = `
		SELECT 
    	id, user_id, token_hash, expires_at, revoked, created_at
		FROM 
			refresh_tokens 
		WHERE 
			id = $1 
		FOR UPDATE`

	var token domain.RefreshToken
	err := r.db.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&token.ID,
		&token.UserID,
		&token.TokenHash,
		&token.ExpiresAt,
		&token.Revoked,
		&token.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrTokenNotFound
		}

		return nil, fmt.Errorf("find token: %w", err)
	}

	return &token, nil
}

func (r *RefreshTokensRepository) Revoked(ctx context.Context, id uuid.UUID) error {
	const query = `UPDATE refresh_tokens SET revoked = true WHERE id = $1 AND revoked = false`

	tag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("update token: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return domain.ErrTokenNotFound
	}

	return nil
}
