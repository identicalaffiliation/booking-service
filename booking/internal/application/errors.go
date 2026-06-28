package application

import "errors"

var (
	ErrRoomAlreadyExists   = errors.New("room already exists")
	ErrInternalRoomUsecase = errors.New("internal server error")

	ErrScheduleAlreadyExists   = errors.New("schedule already exists")
	ErrInvalidRoomId           = errors.New("invalid room id")
	ErrInvalidTimeInterval     = errors.New("invalid start and end time interval")
	ErrInternalScheduleUsecase = errors.New("invalid server error")
)
