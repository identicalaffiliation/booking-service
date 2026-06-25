package integration

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/identicalaffiliation/booking-service/booking/internal/adapters/storage/psql"
	"github.com/identicalaffiliation/booking-service/booking/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSchedulesRepository_CreateSchedule(t *testing.T) {
	schRepo := psql.NewScheduleRepository(db)
	roomsRepo := psql.NewRoomsRepository(db)
	roomID := uuid.New()

	date, err := domain.ParseTimeDate("2026-05-07")
	require.NoError(t, err)

	start, err := domain.ParseTimeDuration("8:00")
	require.NoError(t, err)
	end, err := domain.ParseTimeDuration("18:00")
	require.NoError(t, err)

	actual := domain.NewSchedule(roomID, date, start, end)

	t.Run("success", func(t *testing.T) {
		truncate(db, t)
		_, err := roomsRepo.CreateRoom(context.Background(), &domain.Room{ID: roomID, Name: "test room", Capacity: 5})
		require.NoError(t, err)
		result, err := schRepo.CreateSchedule(context.Background(), actual)
		require.NoError(t, err)

		assert.Equal(t, actual.ID, result.ID)
		assert.Equal(t, actual.RoomID, result.RoomID)
		assert.Equal(t, actual.Day, result.Day)
		assert.Equal(t, actual.EndWorkTime, result.EndWorkTime)
		assert.Equal(t, actual.StartWorkTime, result.StartWorkTime)
	})

	t.Run("error - schedule already exists", func(t *testing.T) {
		truncate(db, t)

		_, err := roomsRepo.CreateRoom(context.Background(), &domain.Room{ID: roomID, Name: "test room", Capacity: 5})
		require.NoError(t, err)
		_, err = schRepo.CreateSchedule(context.Background(), actual)
		require.NoError(t, err)

		result, err := schRepo.CreateSchedule(context.Background(), actual)
		require.Error(t, err)

		assert.ErrorIs(t, err, psql.ErrScheduleAlreadyExists)
		assert.Nil(t, result)
	})
}
