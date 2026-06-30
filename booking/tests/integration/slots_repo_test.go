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

func TestSlotsRepository_CreateSlot(t *testing.T) {
	slotsRepository := psql.NewSlotsRepository(db)
	roomsRepository := psql.NewRoomsRepository(db)

	roomID := uuid.New()
	room := &domain.Room{
		ID:       roomID,
		Name:     "Test Room A",
		Capacity: 5,
	}

	t.Run("success", func(t *testing.T) {
		truncate(db, t)

		_, err := roomsRepository.CreateRoom(context.Background(), room)
		require.NoError(t, err)

		actual := &domain.Slot{
			ID:        uuid.New(),
			RoomID:    roomID,
			Day:       time.Now().UTC().Truncate(time.Hour * 24),
			StartTime: 8 * 60,
			EndTime:   (8 * 60) + (60 * 12),
		}

		err = slotsRepository.CreateSlot(context.Background(), actual)
		require.NoError(t, err)
	})

	t.Run("error - unique violation", func(t *testing.T) {
		truncate(db, t)

		_, err := roomsRepository.CreateRoom(context.Background(), room)
		require.NoError(t, err)

		actual := &domain.Slot{
			ID:        uuid.New(),
			RoomID:    roomID,
			Day:       time.Now().UTC().Truncate(time.Hour * 24),
			StartTime: 8 * 60,
			EndTime:   (8 * 60) + (60 * 12),
		}

		err = slotsRepository.CreateSlot(context.Background(), actual)
		require.NoError(t, err)

		err = slotsRepository.CreateSlot(context.Background(), actual)
		require.Error(t, err)

		assert.ErrorIs(t, err, domain.ErrSlotAlreadyExists)
	})
}
