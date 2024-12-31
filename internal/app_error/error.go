package apperror

type AppError struct {
	ErrorCode int   `json:"error_code"`
	RootError error `json:"root_error"`
}

func (a *AppError) Error() string {
	return a.RootError.Error()
}

var _ error = (*AppError)(nil)
