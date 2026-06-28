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

	actual := domain.NewSchedule(roomID, date, 8*60, (8*60)+(12*60))

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

func TestScheduleRepository_GetAllSchedules(t *testing.T) {
	truncate(db, t)

	roomsRepo := psql.NewRoomsRepository(db)
	scheduleRepository := psql.NewScheduleRepository(db)

	roomID1 := uuid.New()
	date1, err := domain.ParseTimeDate("2026-05-07")
	require.NoError(t, err)

	actual1 := domain.NewSchedule(roomID1, date1, 8*60, (8*60)+(12*60))

	roomID2 := uuid.New()
	date2, err := domain.ParseTimeDate("2026-05-07")
	require.NoError(t, err)

	actual2 := domain.NewSchedule(roomID2, date2, 8*60, (8*60)+(12*60))

	_, err = roomsRepo.CreateRoom(context.Background(), &domain.Room{
		ID:       roomID1,
		Name:     "Test Room A",
		Capacity: 10,
	})
	require.NoError(t, err)

	_, err = roomsRepo.CreateRoom(context.Background(), &domain.Room{
		ID:       roomID2,
		Name:     "Test Room B",
		Capacity: 5,
	})
	require.NoError(t, err)

	schedules := []*domain.Schedule{actual1, actual2}

	for _, schedule := range schedules {
		_, err := scheduleRepository.CreateSchedule(context.Background(), schedule)
		require.NoError(t, err)
	}

	result, err := scheduleRepository.GetAllSchedules(context.Background())
	require.NoError(t, err)

	assert.Equal(t, schedules[0].ID, result[0].ID)
	assert.Equal(t, schedules[0].RoomID, result[0].RoomID)
	assert.Equal(t, schedules[0].Day, result[0].Day)
	assert.Equal(t, schedules[0].StartWorkTime, result[0].StartWorkTime)
	assert.Equal(t, schedules[0].EndWorkTime, result[0].EndWorkTime)

	assert.Equal(t, schedules[1].ID, result[1].ID)
	assert.Equal(t, schedules[1].RoomID, result[1].RoomID)
	assert.Equal(t, schedules[1].Day, result[1].Day)
	assert.Equal(t, schedules[1].StartWorkTime, result[1].StartWorkTime)
	assert.Equal(t, schedules[1].EndWorkTime, result[1].EndWorkTime)
}
