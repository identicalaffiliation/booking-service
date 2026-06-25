package psql

import (
	"context"
	"errors"
	"fmt"

	"github.com/identicalaffiliation/booking-service/booking/internal/domain"
)

var (
	ErrScheduleAlreadyExists = errors.New("schedule already exists")
)

type ScheduleRepository struct {
	db DBTX
}

func NewScheduleRepository(db DBTX) *ScheduleRepository {
	return &ScheduleRepository{db: db}
}

func (r *ScheduleRepository) CreateSchedule(ctx context.Context, s *domain.Schedule) (*domain.Schedule, error) {
	const query = `INSERT INTO schedules (id, room_id, work_day, start_work_time, end_work_time)
		VALUES ($1, $2, $3, $4, $5) RETURNING id, room_id, work_day, start_work_time, end_work_time, created_at`

	var created domain.Schedule
	err := r.db.QueryRow(
		ctx,
		query,
		s.ID,
		s.RoomID,
		s.Day,
		s.StartWorkTime,
		s.EndWorkTime).
		Scan(
			&created.ID,
			&created.RoomID,
			&created.Day,
			&created.StartWorkTime,
			&created.EndWorkTime,
			&created.CreatedAt,
		)
	if err != nil {
		if checkUniqueViolation(err) {
			return nil, ErrScheduleAlreadyExists
		}

		return nil, fmt.Errorf("create schedule: %w", err)
	}

	return &created, nil
}
