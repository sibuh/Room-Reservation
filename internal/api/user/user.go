package user

import (
	"context"
	"errors"
	"net/http"

	"reservation/internal/service/user"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/exp/slog"
)

type UserHandler interface {
	Signup(c *gin.Context)
	Login(c *gin.Context)
	Refresh(c *gin.Context)
}
type userHandler struct {
	logger      *slog.Logger
	userService user.UserService
}

func NewUserHandler(logger *slog.Logger, us user.UserService) UserHandler {
	return &userHandler{
		logger:      logger,
		userService: us,
	}
}

func (uh *userHandler) Signup(c *gin.Context) {
	var req user.SignupRequest
	if err := c.ShouldBind(&req); err != nil {
		uh.logger.Info("failed to bind request input", err.Error())
		c.JSON(http.StatusBadRequest, errors.New("bad request"))
		return
	}
	token, err := uh.userService.Signup(context.Background(), req)
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
func (uh *userHandler) Login(c *gin.Context) {
	var req user.LoginRequest
	if err := c.ShouldBind(&req); err != nil {
		uh.logger.Info("failed to bind login request", err.Error())
		c.JSON(http.StatusBadRequest, err)
		return

	}
	token, err := uh.userService.Login(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
func (uh *userHandler) Refresh(c *gin.Context) {
	userID, ok := c.Value("user_id").(uuid.UUID)
	if !ok {
		uh.logger.Info("user id not set on context", errors.New("user not set on context"))
		c.JSON(http.StatusBadRequest, nil)
		return
	}
	token, err := uh.userService.RefreshToken(context.Background(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})

}
