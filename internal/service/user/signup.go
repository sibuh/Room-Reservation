package user

import (
	"context"
	"errors"
	"reservation/internal/storage/db"
	"reservation/pkg/pass"
	"reservation/pkg/token"
	"time"

	"golang.org/x/exp/slog"
)

type Accesser interface {
	Signup(ctx context.Context, sup SignupRequest) (string, error)
	Login(ctx context.Context, lin LoginRequest) (string, error)
}
type access struct {
	logger slog.Logger
	db.Querier
	key      string
	Duration time.Duration
}

func Init(logger slog.Logger, db db.Querier, key string, dur time.Duration) Accesser {
	return &access{
		logger:   logger,
		Querier:  db,
		key:      key,
		Duration: dur,
	}
}

func (s *access) Signup(ctx context.Context, sup SignupRequest) (string, error) {
	if err := sup.Validate(); err != nil {
		s.logger.Info("invalid input for signup", err.Error())
		return "", err
	}
	hash, err := pass.HashPassword(sup.Password)
	if err != nil {
		s.logger.Error("failed to hash password", err)
		return "", ErrPasswordHash
	}
	sup.Password = hash
	usr, err := s.CreateUser(ctx, db.CreateUserParams{
		FirstName:   sup.FirstName,
		LastName:    sup.LastName,
		Password:    sup.Password,
		PhoneNumber: sup.PhoneNumber,
		Username:    sup.Username,
		Email:       sup.Email,
	})
	if err != nil {
		s.logger.Error("failed to create user", err)
		return "", ErrCreateUser
	}
	// create token
	t, err := token.CreateToken(token.Payload{
		ID:        usr.ID.String(),
		CreatedAt: time.Now(),
		Duration:  s.Duration,
	}, s.key)
	if err != nil {
		s.logger.Error("failed to create token", err)
		return "", err
	}
	return t, nil

}

func (s *access) Login(ctx context.Context, lin LoginRequest) (string, error) {
	if err := lin.Validate(); err != nil {
		s.logger.Info("invalid input for login", err)
		return "", err
	}

	user, err := s.Querier.GetUser(ctx, lin.Email)
	if err != nil {
		return "", err

	}
	if !pass.CheckPasswordHash(lin.Password, user.Password) {
		return "", errors.New("password incorrect")
	}
	tkn, err := token.CreateToken(token.Payload{
		ID:        user.ID.String(),
		CreatedAt: time.Now(),
		Duration:  s.Duration,
	}, s.key)
	if err != nil {
		return "", err
	}
	return tkn, nil
}
