package room

import "errors"

var (
	ErrReservationFailed     = errors.New("failed to reserve room")
	ErrCheckoutSessionFailed = errors.New("failed to create checkout session")
	
)
