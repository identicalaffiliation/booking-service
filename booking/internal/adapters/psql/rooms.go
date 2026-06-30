package psql

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/identicalaffiliation/booking-service/booking/internal/domain"
	"github.com/jackc/pgx/v5"
)

type RoomsRepository struct {
	db DBTX
}

func NewRoomsRepository(db DBTX) *RoomsRepository {
	return &RoomsRepository{
		db: db,
	}
}

func (r *RoomsRepository) CreateRoom(ctx context.Context, room *domain.Room) (*domain.Room, error) {
	const query = `INSERT INTO rooms (id, name, capacity) VALUES ($1, $2, $3)
		RETURNING id, name, capacity, created_at`

	var created domain.Room
	err := r.db.QueryRow(ctx, query, room.ID, room.Name, room.Capacity).Scan(
		&created.ID,
		&created.Name,
		&created.Capacity,
		&created.CreatedAt,
	)
	if err != nil {
		if checkUniqueViolation(err) {
			return nil, domain.ErrRoomAlreadyExists
		}

		return nil, fmt.Errorf("create room: %w", err)
	}

	return &created, nil
}

func (r *RoomsRepository) GetRoom(ctx context.Context, id uuid.UUID) (*domain.Room, error) {
	const query = `SELECT id, name, capacity, created_at FROM rooms WHERE id = $1`

	var room domain.Room
	err := r.db.QueryRow(ctx, query, id).Scan(
		&room.ID,
		&room.Name,
		&room.Capacity,
		&room.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrRoomNotFound
		}

		return nil, fmt.Errorf("find room: %w", err)
	}

	return &room, nil
}

func (r *RoomsRepository) GetRooms(ctx context.Context) ([]*domain.Room, error) {
	const query = `SELECT id, name, capacity, created_at FROM rooms`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("find rooms: %w", err)
	}

	var rooms []*domain.Room
	for rows.Next() {
		var selected domain.Room
		err := rows.Scan(
			&selected.ID,
			&selected.Name,
			&selected.Capacity,
			&selected.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan room: %w", err)
		}

		rooms = append(rooms, &selected)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("scan rooms: %w", err)
	}

	return rooms, nil
}

func (r *RoomsRepository) DeleteRoom(ctx context.Context, id uuid.UUID) error {
	const query = `DELETE FROM rooms WHERE id = $1`

	tag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete room: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return domain.ErrRoomNotFound
	}
	
	return nil
}
