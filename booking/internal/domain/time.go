package domain

import (
	"fmt"
	"time"
)

const (
	TimeLayout = "15:04"
	DateLayout = time.DateOnly
)

func ParseTimeDuration(t string) (int, error) {
	parsed, err := time.Parse(TimeLayout, t)
	if err != nil {
		return 0, fmt.Errorf("parse time: %w", err)
	}

	duration := parsed.Hour()*60 + parsed.Minute()

	return duration, nil
}

func ParseTimeDate(d string) (time.Time, error) {
	parsed, err := time.Parse(DateLayout, d)
	if err != nil {
		return time.Time{}, fmt.Errorf("parse date: %w", err)
	}

	return parsed, nil
}
