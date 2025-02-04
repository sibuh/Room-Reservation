package hotel

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"reservation/internal/apperror"
	"reservation/internal/service/hotel"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

	form, err := c.MultipartForm()
	if err != nil {
		h.logger.Info("unable to get image files from form data", err)
		_ = c.Error(&apperror.AppError{
			ErrorCode: http.StatusBadRequest,
			RootError: errors.New("unable to get image files from form data"),
		})
		return
	}
	values := form.Value
	param.Name = values["name"][0]
	param.City = values["city"][0]
	param.Country = values["country"][0]
	userID, ok := c.Value("user_id").(pgtype.UUID)
	if !ok {
		h.logger.Info("could not get owner id from context", errors.New("failed to get owner id from context"))
		_ = c.Error(&apperror.AppError{
			ErrorCode: http.StatusBadRequest,
			RootError: apperror.ErrInvalidInput,
		})
		return
	}
	param.OwnerID = uuid.MustParse(userID.String())
	param.Rating, err = strconv.ParseFloat(values["rating"][0], 64)
	if err != nil {
		h.logger.Error("failed to parse string to float64 for rating", err)
		_ = c.Error(&apperror.AppError{
			ErrorCode: http.StatusInternalServerError,
			RootError: err,
		})
	}
	if latitude := values["latitude"][0]; latitude != "" {
		param.Location.Latitude, err = strconv.ParseFloat(latitude, 64)
		if err != nil {
			h.logger.Error("failed to parse latitude from multipart form", err)
			_ = c.Error(&apperror.AppError{
				ErrorCode: http.StatusInternalServerError,
				RootError: apperror.ErrBindingQuery,
			})
			return
		}
	}

	if longitude := values["longitude"][0]; longitude != "" {
		param.Location.Longitude, err = strconv.ParseFloat(longitude, 64)
		if err != nil {
			h.logger.Error("failed to parse longitude from multipart form", err)
			_ = c.Error(&apperror.AppError{
				ErrorCode: http.StatusInternalServerError,
				RootError: errors.New("failed to parse longitude from multipart form"),
			})
			return
		}

	}

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
