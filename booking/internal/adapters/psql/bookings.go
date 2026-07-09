package psql

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/identicalaffiliation/booking-service/booking/internal/domain"
	"github.com/jackc/pgx/v5"
)

type BookingsRepository struct {
	db DBTX
}

func NewBookingsRepository(db DBTX) *BookingsRepository {
	return &BookingsRepository{
		db: db,
	}
}

func (r *BookingsRepository) CreateBooking(ctx context.Context, d *domain.Booking) (*domain.Booking, error) {
	const query = `INSERT INTO bookings (id, user_id, slot_id, status)
		VALUES($1, $2, $3, $4) RETURNING id, user_id, slot_id, status, created_at, updated_at`

	var created domain.Booking
	err := r.db.QueryRow(
		ctx,
		query,
		d.ID,
		d.UserID,
		d.SlotID,
		d.Status,
	).Scan(
		&created.ID,
		&created.UserID,
		&created.SlotID,
		&created.Status,
		&created.CreatedAt,
		&created.UpdatedAt,
	)
	if err != nil {
		if checkUniqueViolation(err) {
			return nil, domain.ErrSlotAlreadyBooked
		}

		return nil, fmt.Errorf("create booking: %w", err)
	}

	return &created, nil
}

func (r *BookingsRepository) GetMyBookings(ctx context.Context, userId uuid.UUID) ([]*domain.Booking, error) {
	const query = `SELECT id, user_id, slot_id, status, created_at, updated_at FROM bookings WHERE user_id = $1`

	rows, err := r.db.Query(ctx, query, userId)
	if err != nil {
		return nil, fmt.Errorf("find bookings: %w", err)
	}

	defer rows.Close()

	var bookings []*domain.Booking
	for rows.Next() {
		var row domain.Booking
		err := rows.Scan(
			&row.ID,
			&row.UserID,
			&row.SlotID,
			&row.Status,
			&row.CreatedAt,
			&row.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan booking: %w", err)
		}

		bookings = append(bookings, &row)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("scan bookings: %w", err)
	}

	return bookings, nil
}

func (r *BookingsRepository) GetMyBooking(ctx context.Context, bookingID uuid.UUID) (*domain.MyBooking, error) {
	const query = `
		SELECT
			b.id AS booking_id,
			b.user_id AS user_id,
			b.slot_id AS slot_id,
			s.date AS slot_date,
			s.start_time AS slot_start_time,
			s.end_time AS slot_end_time,
			r.id AS room_id,
			r.name AS room_name,
			r.capacity AS room_capacity,
			b.status AS booking_status,
			b.created_at AS booking_created_at,
			b.updated_at AS booking_updated_at
    FROM 
    	bookings AS b
    INNER JOIN slots AS s ON b.slot_id = s.id
    INNER JOIN rooms AS r ON s.room_id = r.id
    WHERE
    	b.id = $1
	`

	var booking domain.MyBooking
	err := r.db.QueryRow(
		ctx,
		query,
		bookingID,
	).Scan(
		&booking.Booking.ID,
		&booking.Booking.UserID,

		&booking.SlotByBooking.ID,
		&booking.SlotByBooking.Day,
		&booking.SlotByBooking.StartTime,
		&booking.SlotByBooking.EndTime,

		&booking.RoomByBooking.ID,
		&booking.RoomByBooking.Name,
		&booking.RoomByBooking.Capacity,

		&booking.Booking.Status,
		&booking.Booking.CreatedAt,
		&booking.Booking.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrBookingNotFound
		}

		return nil, fmt.Errorf("get booking: %w", err)
	}

	return &booking, nil
}

func (r *BookingsRepository) CancelMyBooking(ctx context.Context, bookingID uuid.UUID) error {
	const query = `UPDATE bookings SET status = 'cancelled', updated_at = now() WHERE id = $1 AND status = 'active'`

	_, err := r.db.Exec(ctx, query, bookingID)
	if err != nil {
		return fmt.Errorf("update booking: %w", err)
	}

	return nil
}
