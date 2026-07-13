package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/identicalaffiliation/booking-service/booking/config"
	"github.com/identicalaffiliation/booking-service/booking/internal/domain"
	"github.com/identicalaffiliation/booking-service/booking/internal/dto/input"
	"github.com/identicalaffiliation/booking-service/booking/internal/dto/output"
	"github.com/identicalaffiliation/booking-service/booking/internal/ports"
)

const (
	userIdKey          = "userId"
	topicCreated       = "bookings.created"
	topicCancelled     = "bookings.cancelled"
	eventTypeCreated   = "created"
	eventTypeCancelled = "cancelled"
)

type BookingsUsecase struct {
	repo   ports.BookingsRepository
	log    ports.Logger
	writer ports.Writer
	cfg    *config.BookingConfig
}

func NewBookingsUsecase(
	repo ports.BookingsRepository,
	log ports.Logger,
	writer ports.Writer,
	cfg *config.BookingConfig,
) *BookingsUsecase {
	return &BookingsUsecase{
		repo:   repo,
		log:    log,
		writer: writer,
		cfg:    cfg,
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

	event := domain.NewEvent(
		created.ID,
		created.UserID,
		created.SlotID,
		eventTypeCreated,
	)

	u.sendEventToBroker(topicCreated, event)

	return output.NewBookingOutput(created), nil
}

func (u *BookingsUsecase) Cancel(ctx context.Context, bookingID uuid.UUID) error {
	if bookingID == uuid.Nil {
		return domain.ErrInvalidBookingData
	}

	if err := u.repo.CancelMyBooking(ctx, bookingID); err != nil {
		return domain.ErrInternal
	}

	booking, err := u.repo.GetMyBooking(ctx, bookingID)
	if err != nil {
		if errors.Is(err, domain.ErrBookingNotFound) {
			return err
		}

		u.log.Error("failed to get booking", "bookingId", bookingID, "error", err)
		return domain.ErrInternal
	}

	event := domain.NewEvent(
		booking.Booking.ID,
		booking.Booking.UserID,
		booking.SlotByBooking.ID,
		eventTypeCancelled,
	)

	u.sendEventToBroker(topicCancelled, event)

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

func (u *BookingsUsecase) sendEventToBroker(topic string, event *domain.KafkaEvent) {
	go func() {
		sendCtx, cancel := context.WithTimeout(context.Background(), u.cfg.Timeout)
		defer cancel()

		if err := u.writer.SendJSON(sendCtx, topic, event.BookingID, event); err != nil {
			u.log.Error("failed to send kafka event", "topic", topic, "error", err)
		}

		u.log.Debug("event is sent")
	}()
}
