package middleware

import (
	"context"
	"errors"
	"net/http"
	"reservation/internal/apperror"
	"reservation/internal/storage/db"
	"reservation/pkg/token"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/exp/slog"
)

const Bearer = "Bearer"

type Middleware interface {
	Authorize() gin.HandlerFunc
}

type middleware struct {
	key    string
	logger *slog.Logger
	db.Querier
}

func NewMiddleware(logger *slog.Logger, q db.Querier, key string) Middleware {
	return &middleware{
		key:     key,
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
			return
		}
		slicedToken := strings.Split(tkn, " ")
		if slicedToken[0] != Bearer {
			a.logger.InfoContext(context.Background(), "token is not of bearer type", slicedToken[0])
			c.AbortWithError(http.StatusUnauthorized, errors.New("token is not of type bearer"))
			return
		}
		payload, err := token.VerifyToken(slicedToken[1], a.key, a.logger)
		if err != nil {
			a.logger.Info("invalid token", err)
			c.AbortWithError(http.StatusUnauthorized, err)
			return
		}
		user, err := a.Querier.GetUserByID(context.Background(), pgtype.UUID{
			Bytes: payload.ID,
			Valid: true,
		})
		if err != nil {
			a.logger.InfoContext(context.Background(), "user does not exist")
			c.AbortWithError(http.StatusUnauthorized, errors.New("user not found"))
			return
		}
		c.Set("user_id", user.ID)

	}
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if err := c.Err(); err != nil {
			thrownError := err.(*apperror.AppError)
			c.JSON(thrownError.ErrorCode, thrownError.RootError)
		}
	}
}
