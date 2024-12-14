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
			a.logger.InfoContext(context.Background(), "no authorization token in request", tkn)
			c.AbortWithError(http.StatusUnauthorized, errors.New("unable to access"))
		}
		slicedToken := strings.Split(tkn, " ")
		if slicedToken[0] != Bearer {
			a.logger.InfoContext(context.Background(), "token is not of bearer type", slicedToken[0])
			c.AbortWithError(http.StatusUnauthorized, errors.New("token is not of type bearer"))
		}
		payload, err := token.VerifyToken(slicedToken[1], a.logger)
		user, err := a.Querier.GetUserByID(context.Background(), uuid.MustParse(payload.ID))
		if err != nil {
			a.logger.InfoContext(context.Background(), "user does not exist")
			c.AbortWithError(http.StatusUnauthorized, errors.New("user not found"))
		}
		c.Set("user_id", user.ID)

	}
}
