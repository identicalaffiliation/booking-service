package psql

import (
	"context"
	"errors"
	"fmt"

	"github.com/identicalaffiliation/booking-service/booking/internal/domain"
)

var (
	ErrSlotAlreadyExists = errors.New("slot already exists")
)

type SlotsRepository struct {
	db DBTX
}

func NewSlotsRepository(db DBTX) *SlotsRepository {
	return &SlotsRepository{
		db: db,
	}
}

func (r *SlotsRepository) CreateSlot(ctx context.Context, slot *domain.Slot) error {
	const query = `INSERT INTO slots (id, room_id, date, start_time, end_time)
		VALUES ($1, $2, $3, $4, $5) RETURNING id, room_id, date, start_time, end_time, created_at`

	_, err := r.db.Exec(ctx, query, slot.ID, slot.RoomID, slot.Day, slot.StartTime, slot.EndTime)
	if err != nil {
		if checkUniqueViolation(err) {
			return ErrSlotAlreadyExists
		}

		return fmt.Errorf("create slot: %w", err)
	}

	return nil
}
