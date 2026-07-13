package psql

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type NotificationsRepository struct {
	db DBTX
}

func NewNotificationsRepository(db DBTX) *NotificationsRepository {
	return &NotificationsRepository{
		db: db,
	}
}

func (r *NotificationsRepository) Create(ctx context.Context, id uuid.UUID, payload []byte) error {
	const query = `INSERT INTO notifications (id, payload) VALUES ($1, $2) ON CONFLICT DO NOTHING`

	_, err := r.db.Exec(ctx, query, id, payload)
	if err != nil {
		return fmt.Errorf("create notification: %w", err)
	}

	return nil
}
