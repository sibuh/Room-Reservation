package hotel

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"reservation/internal/apperror"
	"reservation/internal/storage/db"

	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/exp/slog"
)

type HotelService interface {
	Register(ctx context.Context, param RegisterHotelParam) (db.Hotel, error)
	SearchHotels(ctx context.Context, param SearchHotelParam) ([]db.Hotel, error)
	GetHotels(ctx context.Context) ([]db.Hotel, error)
	GetHotelByName(ctx context.Context, hotelName string) (db.Hotel, error)
}

type hotelService struct {
	db.Querier
	logger *slog.Logger
}

func NewHotelService(q db.Querier, logger *slog.Logger) HotelService {
	return &hotelService{
		Querier: q,
		logger:  logger,
	}
}

func (h *hotelService) Register(ctx context.Context, param RegisterHotelParam) (db.Hotel, error) {
	if err := param.Validate(); err != nil {
		h.logger.Info("invalid input", err)
		return db.Hotel{}, &apperror.AppError{
			ErrorCode: http.StatusBadRequest,
			RootError: apperror.ErrInvalidInput,
		}
	}
	htl, err := h.CreateHotel(ctx, db.CreateHotelParams{
		Name:    param.Name,
		City:    param.City,
		Country: param.Country,
		OwnerID: pgtype.UUID{
			Bytes: param.OwnerID,
			Valid: true,
		},
		Location:  []float64{param.Location.Latitude, param.Location.Longitude},
		Rating:    param.Rating,
		ImageUrls: param.ImageURLs,
	})

	if err != nil {
		h.logger.Error("failed to register hotel", err)
		return db.Hotel{}, &apperror.AppError{
			ErrorCode: http.StatusInternalServerError,
			RootError: apperror.ErrUnableToCreate,
		}
	}

	return htl, nil

}

// TODO:dynamic price calculation must be handled
func (h *hotelService) SearchHotels(ctx context.Context, param SearchHotelParam) ([]db.Hotel, error) {
	data, err := h.Querier.SearchHotels(ctx, db.SearchHotelsParams{
		City:     param.Place,
		Capacity: param.Capacity,
		FromTime: pgtype.Timestamptz{
			Time:  param.FromTime,
			Valid: true,
		},
		FromTime_2: pgtype.Timestamptz{
			Time:  param.ToTime,
			Valid: true,
		},
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &apperror.AppError{
				ErrorCode: http.StatusNotFound,
				RootError: apperror.ErrRecordNotFound}
		}
		return nil, &apperror.AppError{
			ErrorCode: http.StatusInternalServerError,
			RootError: apperror.ErrUnableToGet,
		}
	}
	if len(data) == 0 {
		h.logger.Info("could not get hotel", errors.New("no hotel found for your search"))
		return nil, &apperror.AppError{
			ErrorCode: http.StatusNotFound,
			RootError: errors.New("no hotel found for your search"),
		}
	}
	var hotels []db.Hotel
	for _, v := range data {
		hotels = append(hotels, db.Hotel{
			ID:        v.ID,
			Name:      v.Name,
			Rating:    v.Rating,
			Country:   v.Country,
			City:      v.City,
			Location:  v.Location,
			ImageUrls: v.ImageUrls,
		})
	}
	return hotels, nil
}

func (h *hotelService) GetHotels(ctx context.Context) ([]db.Hotel, error) {
	hotels, err := h.Querier.GetHotels(ctx)
	if err != nil {
		h.logger.Info("failed to get hotels", err)
		return nil, err
	}
	return hotels, nil
}
func (h *hotelService) GetHotelByName(ctx context.Context, hotelName string) (db.Hotel, error) {
	hotel, err := h.Querier.GetHotelByName(ctx, hotelName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.logger.Info("the requested hotel not found",
				fmt.Sprintf("hotelName:%s", hotelName), err)
			return db.Hotel{}, &apperror.AppError{
				ErrorCode: http.StatusNotFound,
				RootError: apperror.ErrRecordNotFound}
		}
		h.logger.Error("unable to get hotel",
			fmt.Sprintf("hotelName:%s", hotelName), err)

		return db.Hotel{}, &apperror.AppError{
			ErrorCode: http.StatusInternalServerError,
			RootError: apperror.ErrUnableToGet}
	}
	return hotel, nil
}
