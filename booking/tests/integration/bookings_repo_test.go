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

func TestBookingsRepository_CreateBooking(t *testing.T) {
	truncate(db, t)

	roomsRepo := psql.NewRoomsRepository(db)
	testRoom := domain.NewRoom("test room", 5)

	_, err := roomsRepo.CreateRoom(context.Background(), testRoom)
	require.NoError(t, err)

	testUser := domain.NewUser("test user", "test password", string(domain.Client))
	testSlot := domain.NewSlot(testRoom.ID, time.Now().UTC(), 900, 1400)
	usersRepo := psql.NewUsersRepository(db)
	slotsRepo := psql.NewSlotsRepository(db)

	_, err = usersRepo.CreateUser(context.Background(), testUser)
	require.NoError(t, err)

	require.NoError(t, slotsRepo.CreateSlot(context.Background(), testSlot))

	t.Run("success", func(t *testing.T) {
		ctx := context.Background()
		repo := psql.NewBookingsRepository(db)

		expected := domain.NewBooking(testUser.ID, testSlot.ID)

		actual, err := repo.CreateBooking(ctx, expected)
		require.NoError(t, err)

		assert.Equal(t, expected.ID, actual.ID)
		assert.Equal(t, expected.UserID, actual.UserID)
		assert.Equal(t, expected.SlotID, actual.SlotID)
		assert.Equal(t, expected.Status, actual.Status)
	})

	t.Run("slot already booked - race condition", func(t *testing.T) {
		ctx := context.Background()

		repo := psql.NewBookingsRepository(db)

		expected := domain.NewBooking(testUser.ID, testSlot.ID)
		actual, err := repo.CreateBooking(ctx, expected)
		require.Error(t, err)

		assert.ErrorIs(t, err, domain.ErrSlotAlreadyBooked)
		assert.Nil(t, actual)
	})
}

func TestBookingsRepository_GetMyBookings(t *testing.T) {
	truncate(db, t)

	roomsRepo := psql.NewRoomsRepository(db)
	testRoom := domain.NewRoom("test room", 5)

	_, err := roomsRepo.CreateRoom(context.Background(), testRoom)
	require.NoError(t, err)

	testUser := domain.NewUser("test user", "test password", string(domain.Client))
	testSlot := domain.NewSlot(testRoom.ID, time.Now().UTC(), 900, 1400)
	usersRepo := psql.NewUsersRepository(db)
	slotsRepo := psql.NewSlotsRepository(db)

	_, err = usersRepo.CreateUser(context.Background(), testUser)
	require.NoError(t, err)

	require.NoError(t, slotsRepo.CreateSlot(context.Background(), testSlot))

	t.Run("success", func(t *testing.T) {
		ctx := context.Background()
		repo := psql.NewBookingsRepository(db)

		expected := domain.NewBooking(testUser.ID, testSlot.ID)
		_, err := repo.CreateBooking(ctx, expected)
		require.NoError(t, err)

		bookings, err := repo.GetMyBookings(ctx, expected.UserID)
		require.NoError(t, err)

		assert.Len(t, bookings, 1)
		assert.Equal(t, expected.ID, bookings[0].ID)
		assert.Equal(t, expected.UserID, bookings[0].UserID)
		assert.Equal(t, expected.SlotID, bookings[0].SlotID)
		assert.Equal(t, expected.Status, bookings[0].Status)
	})
}

func TestBookingsRepository_GetMyBooking(t *testing.T) {
	truncate(db, t)

	roomsRepo := psql.NewRoomsRepository(db)
	testRoom := domain.NewRoom("test room", 5)

	_, err := roomsRepo.CreateRoom(context.Background(), testRoom)
	require.NoError(t, err)

	testUser := domain.NewUser("test user", "test password", string(domain.Client))
	testSlot := domain.NewSlot(testRoom.ID, time.Now().UTC(), 900, 1400)
	usersRepo := psql.NewUsersRepository(db)
	slotsRepo := psql.NewSlotsRepository(db)

	_, err = usersRepo.CreateUser(context.Background(), testUser)
	require.NoError(t, err)

	require.NoError(t, slotsRepo.CreateSlot(context.Background(), testSlot))

	t.Run("success", func(t *testing.T) {
		ctx := context.Background()
		repo := psql.NewBookingsRepository(db)

		expected := domain.NewBooking(testUser.ID, testSlot.ID)

		_, err := repo.CreateBooking(ctx, expected)
		require.NoError(t, err)

		actual, err := repo.GetMyBooking(ctx, expected.ID)
		require.NoError(t, err)

		assert.Equal(t, actual.Booking.ID, expected.ID)
		assert.Equal(t, actual.Booking.UserID, expected.UserID)
		assert.Equal(t, actual.Booking.Status, expected.Status)

		assert.Equal(t, actual.RoomByBooking.ID, testRoom.ID)
		assert.Equal(t, actual.RoomByBooking.Name, testRoom.Name)
		assert.Equal(t, actual.RoomByBooking.Capacity, testRoom.Capacity)

		assert.Equal(t, actual.SlotByBooking.ID, testSlot.ID)
		assert.Equal(t, actual.SlotByBooking.StartTime, testSlot.StartTime)
		assert.Equal(t, actual.SlotByBooking.EndTime, testSlot.EndTime)
	})

	t.Run("booking not found", func(t *testing.T) {
		ctx := context.Background()
		repo := psql.NewBookingsRepository(db)

		actual, err := repo.GetMyBooking(ctx, uuid.New())
		require.Error(t, err)

		assert.ErrorIs(t, err, domain.ErrBookingNotFound)
		assert.Nil(t, actual)
	})
}

func TestBookingsRepository_CancelMyBooking(t *testing.T) {
	truncate(db, t)

	roomsRepo := psql.NewRoomsRepository(db)
	testRoom := domain.NewRoom("test room", 5)

	_, err := roomsRepo.CreateRoom(context.Background(), testRoom)
	require.NoError(t, err)

	testUser := domain.NewUser("test user", "test password", string(domain.Client))
	testSlot := domain.NewSlot(testRoom.ID, time.Now().UTC(), 900, 1400)
	usersRepo := psql.NewUsersRepository(db)
	slotsRepo := psql.NewSlotsRepository(db)

	_, err = usersRepo.CreateUser(context.Background(), testUser)
	require.NoError(t, err)

	require.NoError(t, slotsRepo.CreateSlot(context.Background(), testSlot))

	t.Run("success", func(t *testing.T) {
		ctx := context.Background()
		repo := psql.NewBookingsRepository(db)

		booking := domain.NewBooking(testUser.ID, testSlot.ID)

		_, err := repo.CreateBooking(ctx, booking)
		require.NoError(t, err)

		err = repo.CancelMyBooking(ctx, booking.ID)
		require.NoError(t, err)

		actual, err := repo.GetMyBooking(ctx, booking.ID)
		require.NoError(t, err)
		
		assert.Equal(t, domain.Cancelled, actual.Booking.Status)
	})
}
