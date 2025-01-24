package roomtype

import (
	"context"
	"net/http"
	"reservation/internal/apperror"
	roomtype "reservation/internal/service/room_type"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
)

type RoomTypeHandler interface {
	CreateRoomType(c *gin.Context)
	GetRoomTypes(c *gin.Context)
}

type roomType struct {
	logger *slog.Logger
	svc    roomtype.RoomType
}

func NewRoomTypeHandler(logger *slog.Logger, svc roomtype.RoomType) RoomTypeHandler {
	return &roomType{
		logger: logger,
		svc:    svc,
	}
}
func (rt *roomType) CreateRoomType(c *gin.Context) {
	req := roomtype.CreateRoomTypeRequest{}
	if err := c.ShouldBind(&req); err != nil {
		rt.logger.Info("failed to bind create room type request body", err)
		_ = c.Error(&apperror.AppError{
			ErrorCode: http.StatusBadRequest,
			RootError: apperror.ErrBindingRequestBody,
		})
		return
	}
	res, err := rt.svc.CreateRoomType(context.Background(), req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, res)
}
func (rt *roomType) GetRoomTypes(c *gin.Context) {
	roomTypes, err := rt.svc.GetRoomTypes(context.Background())
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, roomTypes)
}
