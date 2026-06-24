package input

import (
	"errors"
	"strings"
)

var (
	ErrInvalidRoomName     = errors.New("invalid room name")
	ErrInvalidRoomCapacity = errors.New("invalid room capacity")
)

type CreateRoomInput struct {
	Name     string
	Capacity int
}

func NewCreateRoomInput(name string, cap int) *CreateRoomInput {
	return &CreateRoomInput{Name: name, Capacity: cap}
}

func (i *CreateRoomInput) Validate() error {
	if len(strings.TrimSpace(i.Name)) == 0 {
		return ErrInvalidRoomName
	}

	if len(i.Name) < 1 {
		return ErrInvalidRoomCapacity
	}

	if len(i.Name) < 5 {
		return ErrInvalidRoomName
	}

	return nil
}
