package constant

import "errors"

var (
	ErrInvalidBookingTime = errors.New("invalid booking time range")
	ErrSpaceNotFound      = errors.New("space not found")
)
