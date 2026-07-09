package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/identicalaffiliation/booking-service/booking/internal/domain"
	"github.com/identicalaffiliation/booking-service/booking/internal/dto/input"
	"github.com/identicalaffiliation/booking-service/booking/internal/dto/output"
	"github.com/identicalaffiliation/booking-service/booking/internal/ports"
)

type BookingsUsecase struct {
	repo ports.BookingsRepository
	log  ports.Logger
}

func NewBookingsUsecase(repo ports.BookingsRepository, log ports.Logger) *BookingsUsecase {
	return &BookingsUsecase{
		repo: repo,
		log:  log,
	}
}

func (u *BookingsUsecase) Create(
	ctx context.Context,
	in *input.CreateBookingInput,
) (*output.CreateBookingOutput, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	userIdStr := ctx.Value("userId").(string)
	userID, err := uuid.Parse(userIdStr)
	if err != nil {
		return nil, domain.ErrInvalidUserData
	}

	created, err := u.repo.CreateBooking(
		ctx,
		domain.NewBooking(
			userID,
			in.SlotID,
		),
	)
	if err != nil {
		if errors.Is(err, domain.ErrSlotAlreadyBooked) {
			return nil, err
		}

		u.log.Error("failed to booked slot", "error", err)
		return nil, domain.ErrInternal
	}

	return output.NewBookingOutput(created), nil
}

func (u *BookingsUsecase) Cancel(ctx context.Context, bookingID uuid.UUID) error {
	if bookingID == uuid.Nil {
		return domain.ErrInvalidBookingData
	}

	if err := u.repo.CancelMyBooking(ctx, bookingID); err != nil {
		return domain.ErrInternal
	}

	return nil
}
