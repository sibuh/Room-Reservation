package token

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/exp/slog"
)

var ErrExpiredToken = errors.New("token is expired login to get fresh token")
var ErrTokenCreationFailure = errors.New("failed to create token")
var ErrUnexpecteSigningMethod = errors.New("unexpected signing method")

type Payload struct {
	ID        string
	CreatedAt time.Time
	Duration  time.Duration
}

func (p *Payload) Valid() error {
	if time.Now().Before(p.CreatedAt.Add(p.Duration)) {
		return nil
	}
	return ErrExpiredToken
}

func CreateToken(payload Payload, key string, logger slog.Logger) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &payload)

	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		log.Println("failed to create token", err)
		return "", ErrTokenCreationFailure
	}
	return tokenString, nil
}
func VerifyToken(tokenString string, logger slog.Logger) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			log.Println("failed to verify token",
				fmt.Errorf("unexpected signing Method %s", token.Header["alg"]))
			return nil, ErrUnexpecteSigningMethod
		}
		return nil, nil
	}
	token, err := jwt.ParseWithClaims(tokenString, &Payload{}, keyFunc)
	if err != nil {
		log.Println("parsing payload failed")
		return &Payload{}, err
	}
	return token.Claims.(*Payload), nil
}
