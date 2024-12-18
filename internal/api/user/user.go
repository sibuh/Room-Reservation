package user

import (
	"context"
	"errors"
	"net/http"
	usrv "reservation/internal/service/user"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/exp/slog"
)

type User interface {
	Signup(c *gin.Context)
	Login(c *gin.Context)
	Refresh(c *gin.Context)
}
type user struct {
	logger      *slog.Logger
	userService usrv.Accesser
}

func Init(logger *slog.Logger, us usrv.Accesser) User {
	return &user{
		logger:      logger,
		userService: us,
	}
}

func (u *user) Signup(c *gin.Context) {
	var req usrv.SignupRequest
	if err := c.ShouldBind(&req); err != nil {
		u.logger.Info("failed to bind request input", err.Error())
		c.JSON(http.StatusBadRequest, errors.New("bad request"))
		return
	}
	token, err := u.userService.Signup(context.Background(), req)
	if err != nil {
		//TODO: decide the type of error assure if it is internal
		//server error or other like validation error
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})

}
func (u *user) Login(c *gin.Context) {
	var req usrv.LoginRequest
	if err := c.ShouldBind(&req); err != nil {
		u.logger.Info("failed to bind login request", err.Error())
		c.JSON(http.StatusBadRequest, err)
		return

	}
	token, err := u.userService.Login(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
func (u *user) Refresh(c *gin.Context) {
	userID, ok := c.Value("user_id").(uuid.UUID)
	if !ok {
		u.logger.Info("user id not set on context", errors.New("user not set on context"))
		c.JSON(http.StatusBadRequest, nil)
		return
	}
	token, err := u.userService.RefreshToken(context.Background(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})

}
