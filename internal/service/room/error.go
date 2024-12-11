package room

import "errors"

var ErrReservationFailed = errors.New("failed to reserve room")
var ErrCheckoutSessionFailed = errors.New("failed to create checkout session")
var ErrInvalidInput = errors.New("invalid input")
