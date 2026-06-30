package psql

import (
	"context"
	"fmt"

	"github.com/identicalaffiliation/booking-service/booking/internal/domain"
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
