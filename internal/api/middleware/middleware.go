package middleware

import (
	"context"
	"errors"
	"net/http"
	"reservation/internal/storage/db"
	"reservation/pkg/token"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/exp/slog"
)

const Bearer = "Bearer"

var (
	ErrAuthorizedAccess = errors.New("")
)

type Middleware interface {
	Authorize() gin.HandlerFunc
}

type middleware struct {
	logger slog.Logger
	db.Querier
}

func InitMiddleware(logger slog.Logger, q db.Querier) Middleware {
	return &middleware{
		logger:  logger,
		Querier: q,
	}
}

func (a *middleware) Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		tkn := c.Request.Header.Get("Authorization")
		if tkn == "" {
			c.AbortWithError(http.StatusUnauthorized, errors.New("no token not provided"))
		}
		slicedToken := strings.Split(tkn, " ")
		if slicedToken[0] != Bearer {
			c.AbortWithError(http.StatusUnauthorized, errors.New("token is not of type bearer"))
		}
		payload := token.VerifyToken(slicedToken[1])
		user, err := a.Querier.GetUserByID(context.Background(), uuid.MustParse(payload.ID))
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, errors.New("user does not exist"))
		}
		c.Set("user_id", user.ID)

	}
}
