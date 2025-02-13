package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reservation/internal/apperror"
	"reservation/internal/storage/db"
	"reservation/pkg/token"
	"strings"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/exp/slog"
)

const Bearer = "Bearer"

type Middleware interface {
	Authorize() gin.HandlerFunc
	Authenticate() gin.HandlerFunc
}

type middleware struct {
	key    string
	logger *slog.Logger
	db.Querier
	e casbin.IEnforcer
}

func NewMiddleware(logger *slog.Logger, q db.Querier, key string, e casbin.IEnforcer) Middleware {
	return &middleware{
		key:     key,
		logger:  logger,
		Querier: q,
		e:       e,
	}
}

func (m *middleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		tkn := c.Request.Header.Get("Authorization")
		if tkn == "" {
			m.logger.InfoContext(context.Background(), "no authorization token in request", tkn)
			c.AbortWithError(http.StatusUnauthorized, errors.New("unable to access"))
			return
		}
		slicedToken := strings.Split(tkn, " ")
		if slicedToken[0] != Bearer {
			m.logger.InfoContext(context.Background(), "token is not of bearer type", slicedToken[0])
			c.AbortWithError(http.StatusUnauthorized, errors.New("token is not of type bearer"))
			return
		}
		payload, err := token.VerifyToken(slicedToken[1], m.key, m.logger)
		if err != nil {
			m.logger.Info("invalid token", err)
			c.AbortWithError(http.StatusUnauthorized, err)
			return
		}
		user, err := m.Querier.GetUserByID(context.Background(), pgtype.UUID{
			Bytes: payload.ID,
			Valid: true,
		})
		if err != nil {
			m.logger.InfoContext(context.Background(), "user does not exist")
			c.AbortWithError(http.StatusUnauthorized, errors.New("user not found"))
			return
		}
		c.Set("user_id", user.ID)

	}
}
func (m *middleware) Authorize() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sub := ctx.Value("sub")
		dom := ctx.Param("hotel_id")
		obj := ctx.Request.URL
		act := ctx.Request.Method
		b, err := m.e.Enforce(sub, dom, obj, act)
		if err != nil {
			m.logger.Error("failed to enforce authorization policy", err)
			ctx.AbortWithError(http.StatusInternalServerError, errors.New("authorization failed"))
			return
		}
		if !b {
			m.logger.Info("access denied", fmt.Sprintf("sub: %v dom: %v obj: %v act: %v", sub, dom, obj, act))
			ctx.AbortWithError(http.StatusUnauthorized, errors.New("access denied"))
			return
		}
	}
}
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		fmt.Println("time the request took: ",time.Since(start))
		if err := c.Err(); err != nil {
			thrownError := err.(*apperror.AppError)
			c.JSON(thrownError.ErrorCode, thrownError.RootError)
			return
		}
	}
}
