package roomtype

import (
	"context"
	"net/http"
	"reservation/internal/apperror"
	"reservation/internal/storage/db"

	"golang.org/x/exp/slog"
)

type RoomType interface {
	CreateRoomType(ctx context.Context, param CreateRoomTypeRequest) (db.RoomType, error)
}
type roomTypeService struct {
	logger *slog.Logger
	db.Querier
}

func NewRoomTypeService(logger *slog.Logger, q db.Querier) RoomType {
	return &roomTypeService{
		logger:  logger,
		Querier: q,
	}
}

func (rts *roomTypeService) CreateRoomType(ctx context.Context, param CreateRoomTypeRequest) (db.RoomType, error) {
	if err := param.Validate(); err != nil {
		rts.logger.Info("invalid input for create room type request", err)
		return db.RoomType{}, &apperror.AppError{
			ErrorCode: http.StatusBadRequest,
			RootError: apperror.ErrInvalidInput,
		}
	}
	roomType, err := rts.Querier.CreateRoomType(ctx, db.CreateRoomTypeParams{
		RoomType:    db.Roomtype(param.RoomType),
		Description: param.Description,
		Price:       param.Price,
		Capacity:    param.Capacity,
	})
	if err != nil {
		rts.logger.Error("failed to create room type", err)
		return db.RoomType{}, &apperror.AppError{
			ErrorCode: http.StatusInternalServerError,
			RootError: apperror.ErrUnableToCreate,
		}
	}
	return roomType, nil
}
