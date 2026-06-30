package domain

import "errors"

var (
	ErrRoomAlreadyExists = errors.New("room already exists")

	ErrScheduleAlreadyExists = errors.New("schedule already exists")
	ErrInvalidRoomId         = errors.New("invalid room id")
	ErrInvalidTimeInterval   = errors.New("invalid start and end time interval")

	ErrInternal          = errors.New("internal server error")
	ErrSlotAlreadyExists = errors.New("slot already exists")
)
