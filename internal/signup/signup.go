package signup

import (
	"context"
	"reservation/internal/storage/db"
	"reservation/pkg/pass"

	"golang.org/x/exp/slog"
)

type Accesser interface {
	Signup(ctx context.Context, sup SignupRequest) (User, error)
	Login(ctx context.Context, lin LoginRequest) string
}
type signup struct {
	log slog.Logger
	db.Querier
}

func (s *signup) Signup(ctx context.Context, sup SignupRequest) (User, error) {
	if err := sup.Validate(); err != nil {
		s.log.Info(ctx, "invalid input for signup", err.Error())
		return User{}, err
	}
	hash, err := pass.HashPassword(sup.Password)
	if err != nil {
		s.log.Error(ctx, "failed to hash password", err)
		return User{}, ErrPasswordHash
	}
	sup.Password = hash
	u, err := s.CreateUser(ctx, db.CreateUserParams{
		FirstName:   sup.FirstName,
		LastName:    sup.LastName,
		Password:    sup.Password,
		PhoneNumber: sup.PhoneNumber,
		Email:       sup.Email,
	})
	if err != nil {
		s.log.Error(ctx, "failed to create user", err)
		return User{}, ErrCreateUser
	}
	return User(u), ErrInvalidInput

}

func (s *signup) Login(ctx context.Context, lin LoginRequest) (string, error) {
	return "", nil
}
