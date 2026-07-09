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

const (
	userIdKey = "userId"
)

type BookingsUsecase struct {
	repo   ports.BookingsRepository
	log    ports.Logger
	writer ports.Writer
}

func NewBookingsUsecase(repo ports.BookingsRepository, log ports.Logger, writer ports.Writer) *BookingsUsecase {
	return &BookingsUsecase{
		repo:   repo,
		log:    log,
		writer: writer,
	}
}

func (u *BookingsUsecase) Create(
	ctx context.Context,
	in *input.CreateBookingInput,
) (*output.CreateBookingOutput, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	userIdStr := ctx.Value(userIdKey).(string)
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

func (u *BookingsUsecase) GetMyBooking(ctx context.Context, bookingID string) (*output.MyBookingOutput, error) {
	id, err := uuid.Parse(bookingID)
	if err != nil {
		return nil, domain.ErrInvalidBookingData
	}

	booking, err := u.repo.GetMyBooking(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrBookingNotFound) {
			return nil, err
		}

		u.log.Error("failed to get booking by id", "bookingId", id, "error", err)
		return nil, domain.ErrInternal
	}

	out := output.NewMyBookingOutput(
		booking.Booking.ID,
		booking.Booking.UserID,
		booking.SlotByBooking.ID,
		booking.RoomByBooking.ID,
		booking.SlotByBooking.Day.Format(domain.DateLayout),
		minutesToTime(booking.SlotByBooking.StartTime),
		minutesToTime(booking.SlotByBooking.EndTime),
		booking.RoomByBooking.Name,
		booking.RoomByBooking.Capacity,
		string(booking.Booking.Status),
		booking.Booking.CreatedAt,
		booking.Booking.UpdatedAt,
	)

	return out, nil
}

func (u *BookingsUsecase) GetMyBookings(ctx context.Context) (*output.MyBookingsOutput, error) {
	strIdFromCtx := ctx.Value(userIdKey).(string)
	id, err := uuid.Parse(strIdFromCtx)
	if err != nil {
		return nil, domain.ErrInvalidUserData
	}

	bookings, err := u.repo.GetMyBookings(ctx, id)
	if err != nil {
		u.log.Error("failed to get my bookings", "userId", id, "error", err)
		return nil, domain.ErrInternal
	}

	return output.NewMyBookingsOutput(bookings), nil
}
