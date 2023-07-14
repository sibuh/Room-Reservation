package token

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Jwt interface {
	CreateToken(paylad Payload) (string, error)
	VerifyToken(token string) *Payload
}
type Payload struct {
	UserName  string
	CreatedAt time.Time
	Duration  time.Duration
}

func (p *Payload) Valid() error {
	if time.Now().Before(p.CreatedAt.Add(p.Duration)) {
		return nil
	}
	return errors.New("token is expired")
}

type jwtMaker struct {
	SymmetricKey string
}

func NewJwtMaker(key string) Jwt {
	return &jwtMaker{
		SymmetricKey: key,
	}
}
func (j *jwtMaker) CreateToken(payload Payload) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &payload)

	tokenString, err := token.SignedString([]byte(j.SymmetricKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
func (j *jwtMaker) VerifyToken(token string) *Payload {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("wrong signing Method")
		}
		return nil, nil
	}
	Token, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		return &Payload{}
	}

	return Token.Claims.(*Payload)
}
