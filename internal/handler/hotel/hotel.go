package hotel

import (
	"context"
	"errors"
	"net/http"
	"path/filepath"
	"strconv"

	"reservation/internal/apperror"
	"reservation/internal/service/hotel"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/exp/slog"
)

type HotelHandler interface {
	Register(c *gin.Context)
	SearchHotel(c *gin.Context)
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
		_ = c.Error(&apperror.AppError{
			ErrorCode: http.StatusBadRequest,
			RootError: apperror.ErrInvalidInput,
		})
		return
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
	values := form.Value
	param.Name = values["name"][0]
	param.OwnerID = uuid.MustParse(values["owner_id"][0])
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
			h.logger.Error("failed to parse latitude from query param", err)
			_ = c.Error(&apperror.AppError{
				ErrorCode: http.StatusInternalServerError,
				RootError: apperror.ErrBindingQuery,
			})
			return
		}
	}

	if longitude := values["location"][0]; longitude != "" {
		param.Location.Latitude, err = strconv.ParseFloat(longitude, 64)
		h.logger.Error("failed to parse longitude from query param", err)
		_ = c.Error(&apperror.AppError{
			ErrorCode: http.StatusInternalServerError,
			RootError: errors.New("failed to parse locatiion from query param"),
		})
		return

	}

	for _, file := range form.File["images"] {
		savePath := filepath.Join("static", "images", file.Filename)
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
	if rating := c.Query("rating"); rating != "" {
		param.Rating, err = strconv.ParseFloat(rating, 64)
		if err != nil {
			_ = c.Error(&apperror.AppError{
				ErrorCode: http.StatusBadRequest,
				RootError: apperror.ErrBindingQuery,
			})
			return
		}
	}

	if latitude := c.Query("latitude"); latitude != "" {
		param.Location.Latitude, err = strconv.ParseFloat(latitude, 64)
		if err != nil {
			h.logger.Error("failed to parse latitude from query string", err)
			_ = c.Error(&apperror.AppError{
				ErrorCode: http.StatusBadRequest,
				RootError: apperror.ErrBindingQuery,
			})
			return
		}
	}
	if longitude := c.Query("longitude"); longitude != "" {
		param.Location.Longitude, err = strconv.ParseFloat(longitude, 64)
		if err != nil {
			h.logger.Error("failed to parse longitude from query string", err)
			_ = c.Error(&apperror.AppError{
				ErrorCode: http.StatusBadRequest,
				RootError: apperror.ErrBindingQuery,
			})
			return
		}
	}
	htl, err := h.service.SearchHotels(context.Background(), param)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, htl)

}
