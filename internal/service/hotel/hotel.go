package hotel

import (
	"context"
	"reservation/internal/storage/db"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/exp/slog"
)

type HotelService interface {
	Register(ctx context.Context, param RegisterHotelParam) (db.Hotel, error)
}

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type RegisterHotelParam struct {
	Name     string   `json:"name"`
	Location Location `json:"location"`
	Rating   float64  `json:"rating"`
}

func (r RegisterHotelParam) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required.Error("name is required")),
		validation.Field(&r.Location, validation.Required.Error("location is required")),
	)
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
		return db.Hotel{}, err
	}
	//TODO: accept hotel image and upload to storage
	htl, err := h.CreateHotel(ctx, db.CreateHotelParams{
		Name:     param.Name,
		Location: []float64{param.Location.Latitude, param.Location.Longitude},
		Rating: pgtype.Float8{
			Float64: param.Rating,
			Valid:   true,
		},
	})

	if err != nil {
		return db.Hotel{}, err
	}

	return htl, nil

}
