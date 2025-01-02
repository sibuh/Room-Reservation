package room

import "errors"

var (
	ErrReservationFailed     = errors.New("failed to reserve room")
	ErrCheckoutSessionFailed = errors.New("failed to create checkout session")
	ErrInvalidInput          = errors.New("invalid input")
	ErrRecordNotFound        = errors.New("resource not found")
	ErrUnableToGet           = errors.New("unable to get")
)
