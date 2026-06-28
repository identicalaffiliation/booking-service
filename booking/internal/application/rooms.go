package application

import (
	"context"
	"errors"

	"github.com/identicalaffiliation/booking-service/booking/internal/adapters/storage/psql"
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
		if errors.Is(err, psql.ErrRoomAlreadyExists) {
			return nil, ErrRoomAlreadyExists
		}

		u.log.Error("failed to create room", "layer", RoomsLayer, "error", err)
		return nil, ErrInternal
	}

	return output.NewCreateRoomOutput(
			created.ID,
			created.Name,
			created.Capacity,
			created.CreatedAt),
		nil
}
