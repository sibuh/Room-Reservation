package hotel

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"reservation/internal/apperror"
	"reservation/internal/storage/db"

	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/exp/slog"
)

type HotelService interface {
	Register(ctx context.Context, param RegisterHotelParam) (db.Hotel, error)
	SearchHotels(ctx context.Context, param SearchHotelParam) ([]SearchHotelResponse, error)
	GetHotels(ctx context.Context) ([]db.Hotel, error)
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
		Name: param.Name,
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
func (h *hotelService) SearchHotels(ctx context.Context, param SearchHotelParam) ([]SearchHotelResponse, error) {
	hotelsWithRoom, err := h.Querier.SearchHotels(ctx, db.SearchHotelsParams{
		City: param.Place,
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
	var hotelsWithRooms []SearchHotelResponse
	for _, v := range hotelsWithRoom {
		hotelsWithRooms = append(hotelsWithRooms, SearchHotelResponse{
			db.Hotel{
				ID:        v.ID,
				Name:      v.Name,
				Rating:    v.Rating,
				Country:   v.Country,
				City:      v.City,
				Location:  v.Location,
				ImageUrls: v.ImageUrls,
			},
			db.Room{
				ID:         v.ID_2,
				RoomNumber: v.RoomNumber,
				Floor:      v.Floor,
				Status:     v.Status_2,
			},
			db.RoomType{
				ID:          v.ID_3,
				RoomType:    v.RoomType,
				Description: v.Description,
				Price:       v.Price,
			},
		})
	}

	return hotelsWithRooms, nil
}

func (h *hotelService) GetHotels(ctx context.Context) ([]db.Hotel, error) {
	hotels, err := h.Querier.GetHotels(ctx)
	if err != nil {
		h.logger.Info("failed to get hotels", err)
		return nil, err
	}
	return hotels, nil
}
