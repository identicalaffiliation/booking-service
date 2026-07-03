package integration

import (
	"context"
	"testing"

	"github.com/identicalaffiliation/booking-service/booking/internal/adapters/psql"
	"github.com/identicalaffiliation/booking-service/booking/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUsersRepository_CreateUser(t *testing.T) {
	repo := psql.NewUsersRepository(db)

	expected := domain.NewUser("test nickname", "test password", string(domain.Admin))

	t.Run("success", func(t *testing.T) {
		truncate(db, t)

		ctx := context.Background()
		actual, err := repo.CreateUser(ctx, expected)
		require.NoError(t, err)

		assert.Equal(t, expected.ID, actual.ID)
		assert.Equal(t, expected.Nickname, actual.Nickname)
		assert.Equal(t, expected.PasswordHash, actual.PasswordHash)
		assert.Equal(t, expected.Role, actual.Role)
	})

	t.Run("error - unique violation", func(t *testing.T) {
		truncate(db, t)

		ctx := context.Background()
		_, err := repo.CreateUser(ctx, expected)
		require.NoError(t, err)

		_, err = repo.CreateUser(ctx, expected)
		require.Error(t, err)

		assert.ErrorIs(t, err, domain.ErrUserAlreadyExists)
	})
}

func TestUsersRepository_GetUser(t *testing.T) {
	repo := psql.NewUsersRepository(db)
	expected := domain.NewUser("test name", "test password", string(domain.Client))

	t.Run("success", func(t *testing.T) {
		truncate(db, t)

		ctx := context.Background()
		_, err := repo.CreateUser(ctx, expected)
		require.NoError(t, err)

		actual, err := repo.GetUser(ctx, expected.Nickname)
		require.NoError(t, err)

		assert.Equal(t, expected.ID, actual.ID)
		assert.Equal(t, expected.Nickname, actual.Nickname)
		assert.Equal(t, expected.PasswordHash, actual.PasswordHash)
		assert.Equal(t, expected.Role, actual.Role)
	})

	t.Run("error - user not found", func(t *testing.T) {
		truncate(db, t)

		ctx := context.Background()

		actual, err := repo.GetUser(ctx, "invalid nickname))")
		require.Error(t, err)

		assert.ErrorIs(t, err, domain.ErrUserNotFound)
		assert.Nil(t, actual)
	})
}
