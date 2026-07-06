package integration

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/identicalaffiliation/booking-service/booking/internal/adapters/psql"
	"github.com/identicalaffiliation/booking-service/booking/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRefreshTokensRepository_CreateRefreshToken(t *testing.T) {
	truncate(db, t)

	ctx := context.Background()

	u, err := psql.NewUsersRepository(
		db,
	).CreateUser(
		ctx,
		domain.NewUser(
			"test",
			"test",
			"client",
		),
	)
	require.NoError(t, err)

	expected := domain.NewRefreshToken(
		uuid.New(),
		u.ID,
		"test hash",
		false,
		time.Now().Add(time.Hour).Unix(),
	)

	actual, err := psql.NewRefreshTokensRepository(
		db,
	).CreateRefreshToken(
		ctx,
		expected,
	)
	require.NoError(t, err)

	assert.Equal(t, expected.ID, actual.ID)
	assert.Equal(t, expected.UserID, actual.UserID)
	assert.Equal(t, expected.TokenHash, actual.TokenHash)
	assert.Equal(t, expected.Revoked, actual.Revoked)
	assert.Equal(t, expected.ExpiresAt, actual.ExpiresAt)
}

func TestRefreshTokensRepository_GetForUpdate(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		truncate(db, t)
		u, err := psql.NewUsersRepository(
			db,
		).CreateUser(
			ctx,
			domain.NewUser(
				"test",
				"test",
				"client",
			),
		)
		require.NoError(t, err)

		expected := domain.NewRefreshToken(
			uuid.New(),
			u.ID,
			"test hash",
			false,
			time.Now().Add(time.Hour).Unix(),
		)

		_, err = psql.NewRefreshTokensRepository(
			db,
		).CreateRefreshToken(
			ctx,
			expected,
		)
		require.NoError(t, err)

		actual, err := psql.NewRefreshTokensRepository(db).GetForUpdateRefreshToken(ctx, expected.ID)
		require.NoError(t, err)

		assert.Equal(t, expected.ID, actual.ID)
		assert.Equal(t, expected.UserID, actual.UserID)
		assert.Equal(t, expected.TokenHash, actual.TokenHash)
		assert.Equal(t, expected.Revoked, actual.Revoked)
		assert.Equal(t, expected.ExpiresAt, actual.ExpiresAt)
	})

	t.Run("error - token not found", func(t *testing.T) {
		truncate(db, t)

		actual, err := psql.NewRefreshTokensRepository(db).GetForUpdateRefreshToken(ctx, uuid.New())
		require.Error(t, err)

		assert.ErrorIs(t, err, domain.ErrTokenNotFound)
		assert.Nil(t, actual)
	})
}

func TestRefreshTokensRepository_Revoke(t *testing.T) {
	u := domain.NewUser(
		"test nickname",
		"test password",
		"admin",
	)

	token := domain.NewRefreshToken(
		uuid.New(),
		u.ID,
		"test hash",
		false,
		time.Now().Add(time.Hour).Unix(),
	)

	t.Run("success", func(t *testing.T) {
		truncate(db, t)

		_, err := psql.NewUsersRepository(
			db,
		).CreateUser(
			context.Background(),
			u,
		)
		require.NoError(t, err)

		_, err = psql.NewRefreshTokensRepository(
			db,
		).CreateRefreshToken(
			context.Background(),
			token,
		)
		require.NoError(t, err)

		err = psql.NewRefreshTokensRepository(
			db,
		).Revoked(
			context.Background(),
			token.ID,
		)
		require.NoError(t, err)

		actual, err := psql.NewRefreshTokensRepository(
			db,
		).GetForUpdateRefreshToken(
			context.Background(),
			token.ID,
		)
		require.NoError(t, err)

		assert.Equal(t, actual.Revoked, true)

	})

	t.Run("error - token not found", func(t *testing.T) {
		truncate(db, t)
		err := psql.NewRefreshTokensRepository(
			db,
		).Revoked(context.Background(), uuid.New())
		require.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrTokenNotFound)
	})
}
