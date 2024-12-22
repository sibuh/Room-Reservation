package token

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"golang.org/x/exp/slog"
)

var (
	ErrExpiredToken           = errors.New("token is expired login to get fresh token")
	ErrTokenCreationFailure   = errors.New("failed to create token")
	ErrUnexpecteSigningMethod = errors.New("unexpected signing method")
	ErrInvalidToken           = errors.New("invalid token")
)

type Payload struct {
	ID        uuid.UUID
	CreatedAt time.Time
	Duration  time.Duration
}

func (p *Payload) Valid() error {
	if time.Now().Before(p.CreatedAt.Add(p.Duration)) {
		return nil
	}
	return ErrExpiredToken
}

func CreateToken(payload Payload, key string, logger *slog.Logger) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &payload)
	if len(key) != 32 {
		logger.Error("invalid signing key length", len(key))
		return "", errors.New("invalid signing key length")
	}
	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		logger.Error("failed to create token", err)
		return "", ErrTokenCreationFailure
	}
	return tokenString, nil
}
func VerifyToken(tokenString, key string, logger *slog.Logger) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(key), nil
	}

	jwtToken, err := jwt.ParseWithClaims(tokenString, &Payload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}
