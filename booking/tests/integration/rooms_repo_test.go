package integration

import (
	"context"
	"testing"
	"time"

	"github.com/identicalaffiliation/booking-service/booking/internal/adapters/psql"
	"github.com/identicalaffiliation/booking-service/booking/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRoomsRepository_CreateRoom(t *testing.T) {
	truncate(db, t)

	repo := psql.NewRoomsRepository(db)
	expected := domain.NewRoom("test room", 10)

	t.Run("success", func(t *testing.T) {
		ctx := context.Background()
		now := time.Now()
		delta := time.Second * 2

		result, err := repo.CreateRoom(ctx, expected)
		require.NoError(t, err)

		assert.Equal(t, expected.ID, result.ID)
		assert.Equal(t, expected.Name, result.Name)
		assert.Equal(t, expected.Capacity, result.Capacity)
		assert.WithinDuration(t, now, result.CreatedAt, delta)
	})

	t.Run("unique violation error", func(t *testing.T) {
		ctx := context.Background()

		result, err := repo.CreateRoom(ctx, expected)
		require.Error(t, err)

		assert.ErrorIs(t, err, domain.ErrRoomAlreadyExists)
		assert.Nil(t, result)
	})
}
