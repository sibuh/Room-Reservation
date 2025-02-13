package hotel

import (
	"context"
	"errors"
	"net/http"

	"reservation/internal/apperror"
	"reservation/internal/service/hotel"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/exp/slog"
)

type HotelHandler interface {
	Register(c *gin.Context)
	SearchHotel(c *gin.Context)
	GetHotels(c *gin.Context)
	GetHotelByName(c *gin.Context)
	VerifyHotel(c *gin.Context)
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
		h.logger.Error("invalid input", err)
		_ = c.Error(&apperror.AppError{
			ErrorCode: http.StatusBadRequest,
			RootError: apperror.ErrInvalidInput,
		})
	}

	form, err := c.MultipartForm()
	if err != nil {
		h.logger.Info("unable to get image files from form data", err)
		_ = c.Error(&apperror.AppError{
			ErrorCode: http.StatusBadRequest,
			RootError: errors.New("unable to get image files from form data"),
		})
		return
	}

	userID, ok := c.Value("user_id").(pgtype.UUID)
	if !ok {
		h.logger.Info("could not get owner id from context", errors.New("failed to get owner id from context"))
		_ = c.Error(&apperror.AppError{
			ErrorCode: http.StatusBadRequest,
			RootError: apperror.ErrInvalidInput,
		})
		return
	}
	param.OwnerID = userID.Bytes

	for _, file := range form.File["images"] {
		savePath := "public/" + file.Filename
		if err := c.SaveUploadedFile(file, savePath); err != nil {
			h.logger.Error("failed to save uploaded file", err)
			_ = c.Error(&apperror.AppError{
				ErrorCode: http.StatusInternalServerError,
				RootError: errors.New("failed to save uploaded file"),
			})

			return
		}
		param.ImageURLs = append(param.ImageURLs, savePath)

	}

	htl, err := h.service.Register(context.Background(), param)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, htl)

}
func (h *hotelHandler) SearchHotel(c *gin.Context) {
	var param hotel.SearchHotelParam
	var err error
	if err := c.ShouldBind(&param); err != nil {
		h.logger.Info("failed to bind search hotel params", err)
		_ = c.Error(&apperror.AppError{
			ErrorCode: http.StatusBadRequest,
			RootError: apperror.ErrBindingRequestBody,
		})
		return
	}
	htl, err := h.service.SearchHotels(context.Background(), param)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, htl)

}
func (h *hotelHandler) GetHotels(c *gin.Context) {
	hotels, err := h.service.GetHotels(context.Background())
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, hotels)

}
func (h *hotelHandler) GetHotelByName(c *gin.Context) {
	hotelName := c.Query("hotel_name")
	hotel, err := h.service.GetHotelByName(context.Background(), hotelName)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, hotel)
}
func (h *hotelHandler) VerifyHotel(c *gin.Context) {
	hotelID := c.Param("hotel_id")
	hotel, err := h.service.VerifyHotel(context.Background(), hotelID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, hotel)
}
