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

const RoomsLayer = "rooms service layer"

type RoomsUsecase struct {
	repo ports.RoomsRepository
	log  ports.Logger
}

func NewRoomsUsecase(repo ports.RoomsRepository, log ports.Logger) *RoomsUsecase {
	return &RoomsUsecase{repo: repo, log: log}
}

func (u *RoomsUsecase) CreateRoom(
	ctx context.Context,
	input *input.CreateRoomInput,
) (*output.CreateRoomOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	room := domain.NewRoom(input.Name, input.Capacity)

	created, err := u.repo.CreateRoom(ctx, room)
	if err != nil {
		if errors.Is(err, domain.ErrRoomAlreadyExists) {
			return nil, err
		}

		u.log.Error("failed to create room", "layer", RoomsLayer, "error", err)
		return nil, domain.ErrInternal
	}

	return output.NewCreateRoomOutput(
			created.ID,
			created.Name,
			created.Capacity,
			created.CreatedAt),
		nil
}

func (u *RoomsUsecase) GetRoom(ctx context.Context, id string) (*output.RoomOutput, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, domain.ErrInvalidRoomData
	}

	room, err := u.repo.GetRoom(ctx, uid)
	if err != nil {
		if errors.Is(err, domain.ErrRoomNotFound) {
			return nil, err
		}

		u.log.Error("failed to find room by id", "room id", uid, "error", err)
		return nil, domain.ErrInternal
	}

	return output.NewRoomOutput(room.ID, room.Name, room.Capacity, room.CreatedAt), nil
}

func (u *RoomsUsecase) GetRooms(ctx context.Context) (*output.RoomsOutput, error) {
	rooms, err := u.repo.GetRooms(ctx)
	if err != nil {
		u.log.Error("failed to find rooms", "error", err)
		return nil, domain.ErrInternal
	}

	return u.roomsDomainToRoomsOutput(rooms), nil
}

func (u *RoomsUsecase) DeleteRoom(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return domain.ErrInvalidRoomData
	}

	if err := u.repo.DeleteRoom(ctx, uid); err != nil {
		if errors.Is(err, domain.ErrRoomNotFound) {
			return err
		}

		u.log.Error("failed to delete room by id", "room id", uid, "error", err)
		return domain.ErrInternal
	}

	return nil
}

func (u *RoomsUsecase) roomsDomainToRoomsOutput(d []*domain.Room) *output.RoomsOutput {
	rooms := make([]*output.RoomOutput, 0, len(d))
	for _, room := range d {
		rooms = append(
			rooms,
			output.NewRoomOutput(
				room.ID,
				room.Name,
				room.Capacity,
				room.CreatedAt,
			),
		)
	}

	return &output.RoomsOutput{Rooms: rooms}
}
