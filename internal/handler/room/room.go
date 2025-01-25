package room

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"reservation/internal/apperror"
	"reservation/internal/service/room"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"golang.org/x/exp/slog"
)

type RoomHandler interface {
	Reserve(c *gin.Context)
	AddRoom(c *gin.Context)
	UpdateRoom(c *gin.Context)
	GetHotelRooms(c *gin.Context)
	GetRoomReservations(c *gin.Context)
}

type roomHandler struct {
	logger         *slog.Logger
	srv            room.RoomService
	publishableKey string
}

func NewRoomHandler(logger *slog.Logger, srv room.RoomService) RoomHandler {
	return &roomHandler{
		logger: logger,
		srv:    srv,
	}
}

func (r *roomHandler) Reserve(c *gin.Context) {
	req := room.ReserveRoom{}
	if err := c.ShouldBind(&req); err != nil {
		r.logger.Info("failed to bind request body", err)
		_ = c.Error(&apperror.AppError{
			ErrorCode: http.StatusBadRequest,
			RootError: apperror.ErrInvalidInput,
		})
		return
	}
	req.UserID = c.Value("user_id").(uuid.UUID)
	url, err := r.srv.ReserveRoom(context.Background(), req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, url)
}
func (r *roomHandler) UpdateRoom(c *gin.Context) {
	req := room.UpdateRoom{}
	if err := c.ShouldBind(&req); err != nil {
		r.logger.Info("failed updated room params", err)
		_ = c.Error(&apperror.AppError{
			ErrorCode: http.StatusBadRequest,
			RootError: apperror.ErrBindingRequestBody,
		})
		return
	}
	room, err := r.srv.UpdateRoom(context.Background(), req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, room)
}

func (rh *roomHandler) GetRoomReservations(c *gin.Context) {
	roomID := c.Param("room_id")
	rvns, err := rh.srv.GetRoomReservations(context.Background(), roomID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			rh.logger.Info("no reservations related to room found", err, fmt.Sprintf("room_id: %s", roomID))
			_ = c.Error(&apperror.AppError{
				ErrorCode: http.StatusNotFound,
				RootError: apperror.ErrRecordNotFound,
			})
			return
		}
		rh.logger.Error("failed to get reservations related room", fmt.Sprintf("room_id: %s", roomID), err)
		_ = c.Error(&apperror.AppError{
			ErrorCode: http.StatusInternalServerError,
			RootError: apperror.ErrUnableToGet,
		})
		return
	}
	c.JSON(http.StatusOK, rvns)
}
func (rh *roomHandler) AddRoom(c *gin.Context) {
	var param room.CreateRoomParam

	if err := c.ShouldBind(&param); err != nil {
		rh.logger.Info("failed to bind the request", err)
		_ = c.Error(&apperror.AppError{
			ErrorCode: http.StatusBadRequest,
			RootError: apperror.ErrInvalidInput,
		})
		return
	}

	res, err := rh.srv.AddRoom(context.Background(), param)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, res)

}
func (rh *roomHandler) GetHotelRooms(c *gin.Context) {
	hotelID := c.Param("hotel_id")
	if hotelID == "" {
		rh.logger.Info("failed to get hotel id from param", errors.New("hotel id is empty"))
		_ = c.Error(&apperror.AppError{
			ErrorCode: http.StatusBadRequest,
			RootError: apperror.ErrInvalidInput,
		})
		return
	}
	res, err := rh.srv.GetHotelRooms(context.Background(), hotelID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, res)

}
