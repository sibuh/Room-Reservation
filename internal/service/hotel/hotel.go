package hotel

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"reservation/internal/apperror"
	"reservation/internal/storage/db"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/exp/slog"
)

type HotelService interface {
	Register(ctx context.Context, param RegisterHotelParam) (db.Hotel, error)
	SearchHotels(ctx context.Context, param SearchHotelParam) ([]db.Hotel, error)
	GetHotels(ctx context.Context) ([]db.Hotel, error)
	GetHotelByName(ctx context.Context, hotelName string) (db.Hotel, error)
	VerifyHotel(ctx context.Context, ID string) (db.Hotel, error)
}

type hotelService struct {
	db.Querier
	*pgxpool.Pool
	logger *slog.Logger
}

func NewHotelService(q db.Querier, logger *slog.Logger, pool *pgxpool.Pool) HotelService {
	return &hotelService{
		Querier: q,
		logger:  logger,
		Pool:    pool,
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
		ToTime: pgtype.Timestamptz{
			Time:  param.FromTime,
			Valid: true,
		}})

	if err != nil {
		h.logger.Error("failed to search hotels", err)
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
	//return least room price
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

func (h *hotelService) VerifyHotel(ctx context.Context, ID string) (db.Hotel, error) {
	hotel, err := h.Querier.VerifyHotel(ctx, pgtype.UUID{
		Bytes: uuid.MustParse(ID),
		Valid: true,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.logger.Info("hotel to be verified not found", err)
			return db.Hotel{}, &apperror.AppError{
				ErrorCode: http.StatusNotFound,
				RootError: apperror.ErrRecordNotFound,
			}
		}
		h.logger.Error("failed to get hotel to be verified", err)
		return db.Hotel{}, &apperror.AppError{
			ErrorCode: http.StatusInternalServerError,
			RootError: apperror.ErrUnableToGet,
		}
	}
	return hotel, nil
}
