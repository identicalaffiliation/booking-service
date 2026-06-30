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

func TestRoomsRepository_GetRoom(t *testing.T) {
	truncate(db, t)

	repo := psql.NewRoomsRepository(db)
	expected := domain.NewRoom("test Room", 5)
	_, err := repo.CreateRoom(context.Background(), expected)
	require.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		ctx := context.Background()
		now := time.Now().UTC()
		delta := time.Minute * 2

		result, err := repo.GetRoom(ctx, expected.ID)
		require.NoError(t, err)

		assert.Equal(t, expected.ID, result.ID)
		assert.Equal(t, expected.Name, result.Name)
		assert.Equal(t, expected.Capacity, result.Capacity)
		assert.WithinDuration(t, now, result.CreatedAt, delta)
	})

	t.Run("error - room not found", func(t *testing.T) {
		ctx := context.Background()

		result, err := repo.GetRoom(ctx, uuid.New())
		require.Error(t, err)

		assert.ErrorIs(t, err, domain.ErrRoomNotFound)
		assert.Nil(t, result)
	})
}
func TestRoomsRepository_GetRooms(t *testing.T) {
	truncate(db, t)

	ctx := context.Background()
	repo := psql.NewRoomsRepository(db)

	room1 := domain.NewRoom("room first", 5)
	room2 := domain.NewRoom("room second", 3)

	_, err := repo.CreateRoom(ctx, room2)
	require.NoError(t, err)
	_, err = repo.CreateRoom(ctx, room1)
	require.NoError(t, err)

	rooms, err := repo.GetRooms(ctx)
	require.NoError(t, err)

	assert.Equal(t, room2.ID, rooms[0].ID)
	assert.Equal(t, room2.Capacity, rooms[0].Capacity)
	assert.Equal(t, room2.Name, rooms[0].Name)

	assert.Equal(t, room1.ID, rooms[1].ID)
	assert.Equal(t, room1.Capacity, rooms[1].Capacity)
	assert.Equal(t, room1.Name, rooms[1].Name)
}

func TestRoomsRepository_DeleteRoom(t *testing.T) {
	truncate(db, t)

	repo := psql.NewRoomsRepository(db)

	room := domain.NewRoom(
		"first room",
		5,
	)
	_, err := repo.CreateRoom(
		context.Background(),
		room,
	)
	require.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		err := repo.DeleteRoom(context.Background(), room.ID)
		require.NoError(t, err)

		rooms, err := repo.GetRooms(context.Background())
		require.NoError(t, err)

		assert.Len(t, rooms, 0)
	})

	t.Run("error - room not found", func(t *testing.T) {
		err := repo.DeleteRoom(context.Background(), uuid.New())
		require.Error(t, err)
		
		assert.ErrorIs(t, err, domain.ErrRoomNotFound)
	})
}
