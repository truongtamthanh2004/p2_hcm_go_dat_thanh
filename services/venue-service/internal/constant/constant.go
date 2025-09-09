package constant

import "errors"

var (
	ErrUnauthorized        = errors.New("unauthorized")
	ErrInvalidID           = errors.New("invalid id")
	ErrNotFound            = errors.New("not found")
	ErrBadRequest          = errors.New("bad request")
	ErrInternal            = errors.New("internal server error")
	ErrCreateFailed        = errors.New("failed to create resource")
	ErrUpdateFailed        = errors.New("failed to update resource")
	ErrDeleteFailed        = errors.New("failed to delete resource")
	ErrForbidden           = errors.New("forbidden")
	ErrAlreadyExists       = errors.New("amenity already added")
	ErrInternalServerError = errors.New("internal server error")
	ErrAmenityNotFound     = errors.New("amenity not found")
	ErrVenueNotFound       = errors.New("venue not found")
	ErrInvalidSpaceType    = errors.New("invalid space type")
)

const (
	PENDING  = "pending"
	APPROVED = "approved"
	BLOCKED  = "blocked"
)

const (
	PRIVATE_OFFICE = "private_office"
	MEETING_ROOM   = "meeting_room"
	DESK           = "desk"
)
