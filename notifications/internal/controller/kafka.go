package controller

import (
	"context"

	"github.com/google/uuid"
	"github.com/identicalaffiliation/booking-service/notifications/internal/ports"
)

type BookingsHandler struct {
	repo ports.Repository
}

func NewBookingsHandler(repo ports.Repository) *BookingsHandler {
	return &BookingsHandler{repo: repo}
}

func (h *BookingsHandler) HandleBookingMessage(ctx context.Context, value []byte) error {
	err := h.repo.Create(ctx, uuid.New(), value)
	if err != nil {
		return err
	}

	return nil
}
