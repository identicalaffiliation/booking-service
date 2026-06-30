package domain

import "errors"

var (
	ErrRoomAlreadyExists     = errors.New("room already exists")
	ErrScheduleAlreadyExists = errors.New("schedule already exists")
	ErrSlotAlreadyExists     = errors.New("slot already exists")

	ErrInvalidRoomData     = errors.New("invalid room data")
	ErrInvalidScheduleData = errors.New("invalid schedule data")

	ErrInternal = errors.New("internal server error")
)
