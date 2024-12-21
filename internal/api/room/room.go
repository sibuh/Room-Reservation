package room

import (
	"context"
	"net/http"

	"reservation/internal/service/room"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/exp/slog"
)

type RoomHandler interface {
	Reserve(c *gin.Context)
}

type roomHandler struct {
	logger *slog.Logger
	srv    room.RoomService
}

func NewRoomHandler(logger *slog.Logger, srv room.RoomService)

func (r *roomHandler) Reserve(c *gin.Context) {
	req := room.ReserveRoom{}
	if err := c.ShouldBind(&req); err != nil {
		r.logger.Info("failed to bind request body", err)
		c.JSON(http.StatusBadRequest, err)
		return
	}
	req.UserID = c.Value("user_id").(uuid.UUID)
	url, err := r.srv.ReserveRoom(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, url)
}
func (r *roomHandler) UpdateRoom(c *gin.Context) {
	req := room.UpdateRoom{}
	if err := c.ShouldBind(&req); err != nil {
		r.logger.Info("failed updated room params", err)
		c.JSON(http.StatusBadRequest, err)
		return
	}
	room, err := r.srv.UpdateRoom(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, room)
}
