package signup

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
	key string
}

func Init(logger slog.Logger, db db.Querier, key string) Accesser {
	return &access{
		logger:  logger,
		Querier: db,
		key:     key,
	}
}

func (s *access) Signup(ctx context.Context, sup SignupRequest) (string, error) {
	if err := sup.Validate(); err != nil {
		s.logger.Info(ctx, "invalid input for signup", err.Error())
		return "", err
	}
	hash, err := pass.HashPassword(sup.Password)
	if err != nil {
		s.logger.Error(ctx, "failed to hash password", err)
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
		s.logger.Error(ctx, "failed to create user", err)
		return "", ErrCreateUser
	}
	// create token
	t, err := token.CreateToken(token.Payload{
		ID:        usr.ID.String(),
		CreatedAt: time.Now(),
		Duration:  5 * time.Minute, //this should come from config
	}, s.key)
	if err != nil {
		s.logger.Error(ctx, "failed to create token", err)
		return "", err
	}
	return t, nil

}

func (s *access) Login(ctx context.Context, lin LoginRequest) (string, error) {
	if err := lin.Validate(); err != nil {
		s.logger.Info(ctx, "invalid input for login", err)
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
		Duration:  5 * time.Minute,
	}, s.key)
	if err != nil {
		return "", err
	}
	return tkn, nil
}
