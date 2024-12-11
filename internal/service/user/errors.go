package user

import "errors"

var ErrInvalidInput = errors.New("invalid input during sign up")
var ErrPasswordHash = errors.New("failed to hash password")
var ErrCreateUser = errors.New("failed to create user")
