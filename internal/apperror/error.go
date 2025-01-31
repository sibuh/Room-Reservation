package apperror

import "errors"

var (
	ErrInvalidInput       = errors.New("invalid input")
	ErrRecordNotFound     = errors.New("resource not found")
	ErrUnableToGet        = errors.New("unable to get")
	ErrUnableToCreate     = errors.New("unable to create")
	ErrBindingQuery       = errors.New("failed to bind query param")
	ErrBindingRequestBody = errors.New("failed to bind request body")
	ErrCapturingPayment   = errors.New("failed to capture paypal order payment")
)

type AppError struct {
	ErrorCode int   `json:"error_code"`
	RootError error `json:"root_error"`
}

func (a *AppError) Error() string {
	return a.RootError.Error()
}

var _ error = (*AppError)(nil)
