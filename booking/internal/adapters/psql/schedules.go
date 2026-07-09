package psql

import (
	"context"
	"fmt"
	"time"

	"github.com/identicalaffiliation/booking-service/booking/internal/domain"
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
			return nil, domain.ErrScheduleAlreadyExists
		}

		return nil, fmt.Errorf("create schedule: %w", err)
	}

	return &created, nil
}

func (r *ScheduleRepository) GetAllSchedulesByToday(ctx context.Context, date time.Time) ([]*domain.Schedule, error) {
	const query = `SELECT id, room_id, work_day, start_work_time, end_work_time, created_at 
		FROM schedules WHERE work_day = $1`

	rows, err := r.db.Query(ctx, query, date)
	if err != nil {
		return nil, fmt.Errorf("get schedules: %w", err)
	}

	defer rows.Close()

	var schedules []*domain.Schedule
	for rows.Next() {
		var schedule domain.Schedule
		err := rows.Scan(
			&schedule.ID,
			&schedule.RoomID,
			&schedule.Day,
			&schedule.StartWorkTime,
			&schedule.EndWorkTime,
			&schedule.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan schedule: %w", err)
		}

		schedules = append(schedules, &schedule)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("after scan: %w", err)
	}

	return schedules, nil
}
