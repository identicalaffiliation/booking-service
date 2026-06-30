package psql

import (
	"context"
	"fmt"

	"github.com/identicalaffiliation/booking-service/booking/internal/domain"
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
		VALUES ($1, $2, $3, $4, $5)`

	_, err := r.db.Exec(ctx, query, slot.ID, slot.RoomID, slot.Day, slot.StartTime, slot.EndTime)
	if err != nil {
		if checkUniqueViolation(err) {
			return domain.ErrSlotAlreadyExists
		}

		return fmt.Errorf("create slot: %w", err)
	}

	return nil
}
