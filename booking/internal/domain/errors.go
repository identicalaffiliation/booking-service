package domain

import "errors"

var (
	ErrRoomNotFound     = errors.New("room not found")
	ErrScheduleNotFound = errors.New("schedule not found")
	ErrUserNotFound     = errors.New("user not found")
	ErrTokenNotFound    = errors.New("token not found")

	ErrRoomAlreadyExists     = errors.New("room already exists")
	ErrScheduleAlreadyExists = errors.New("schedule already exists")
	ErrSlotAlreadyExists     = errors.New("slot already exists")
	ErrUserAlreadyExists     = errors.New("user already exists")

	ErrInvalidRoomData         = errors.New("invalid room data")
	ErrInvalidScheduleData     = errors.New("invalid schedule data")
	ErrInvalidUserData         = errors.New("invalid user data")
	ErrInvalidRefreshTokenData = errors.New("invalid refresh token data")

	ErrInternal = errors.New("internal server error")
)
