package hotel

import (
	"context"
	"fmt"
	"net/http"
	"reservation/internal/apperror"
	"reservation/internal/storage/db"

	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/exp/slog"
)

type HotelService interface {
	Register(ctx context.Context, param RegisterHotelParam) (db.Hotel, error)
	SearchHotels(ctx context.Context, param SearchHotelParam) (db.Hotel, error)
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
func (h *hotelService) SearchHotels(ctx context.Context, param SearchHotelParam) (db.Hotel, error) {

	return db.Hotel{}, nil
}

func (h *hotelService) GetHotels(ctx context.Context) ([]db.Hotel, error) {
	hotels, err := h.Querier.GetHotels(ctx)
	if err != nil {
		h.logger.Info("failed to get hotels", err)
		return nil, err
	}
	return hotels, nil
}
func BuildSearchQuery(tableName string, param SearchHotelParam) string {
	query := `
        SELECT *
        FROM hotels 
        WHERE TRUE
    `
	var condition string
	if param.City != "" {
		condition = fmt.Sprintf("AND city =%s", param.City)
	}
	query = query + condition
	return query

}
