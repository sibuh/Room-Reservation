package hotel

import (
	"context"
	"net/http"

	"reservation/internal/service/hotel"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
)

type HotelHandler interface {
	Register(c *gin.Context)
}
type hotelHandler struct {
	logger  *slog.Logger
	service hotel.HotelService
}

func NewHotelHandler(logger *slog.Logger, svc hotel.HotelService) HotelHandler {
	return &hotelHandler{
		logger:  logger,
		service: svc,
	}
}
func (h *hotelHandler) Register(c *gin.Context) {
	param := hotel.RegisterHotelParam{}

	if err := c.ShouldBind(&param); err != nil {
		h.logger.Info("failed to bind request body", err)
		c.JSON(http.StatusOK, err)
		return
	}

	htl, err := h.service.Register(context.Background(), param)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, htl)

}
