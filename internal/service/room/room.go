package room

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"reservation/internal/apperror"
	"reservation/internal/storage/db"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"golang.org/x/exp/slog"
)

type RoomService interface {
	ReserveRoom(ctx context.Context, param ReserveRoom) (db.Reservation, error)
	UpdateRoom(ctx context.Context, param UpdateRoom) (Room, error)
	GetRoomReservations(ctx context.Context, roomID string) ([]db.Reservation, error)
	SearchRoom(ctx context.Context, searchParam SearchParam) ([]db.Room, error)
	AddRoom(ctx context.Context, param CreateRoomParam) (CreateRoomResponse, error)
}

type ReservationStatus string

const (
	StatusPending    ReservationStatus = "PENDING"
	StatusSuccessful ReservationStatus = "SUCCESSFUL"
	StatusFailed     ReservationStatus = "FAILED"
)

type roomService struct {
	db.Querier
	*pgxpool.Pool
	logger          *slog.Logger
	stripeSecretKey string
}

func NewRoomService(pool *pgxpool.Pool, q db.Querier, logger *slog.Logger, key string) RoomService {
	return &roomService{
		Querier:         q,
		stripeSecretKey: key,
		logger:          logger,
		Pool:            pool,
	}
}

func (rs *roomService) ReserveRoom(ctx context.Context, param ReserveRoom) (db.Reservation, error) {
	if err := param.Validate(); err != nil {

		return db.Reservation{}, &apperror.AppError{
			ErrorCode: http.StatusBadRequest,
			RootError: err,
		}
	}

	rvn, err := rs.CreateReservation(ctx, db.CreateReservationParams{
		RoomID:      pgtype.UUID{Bytes: param.RoomID, Valid: true},
		FirstName:   param.FirstName,
		LastName:    param.LastName,
		PhoneNumber: param.PhoneNumber,
		Email:       param.Email,
		Status:      db.ReservationStatus(StatusPending),
		FromTime:    pgtype.Timestamptz{Time: param.FromTime, Valid: true},
		ToTime:      pgtype.Timestamptz{Time: param.ToTime, Valid: true},
	})
	if err != nil {
		rs.logger.Error("failed to create reservation", err)
		return db.Reservation{}, &apperror.AppError{
			ErrorCode: http.StatusInternalServerError,
			RootError: errors.New("failed to make reservation"),
		}
	}
	// secretKey, err := rs.createPaymentIntent(ctx, rvn.ID.String(), param.RoomID.String())
	// if err != nil {
	// 	return "", err
	// }
	return rvn, nil

}
func (rs *roomService) UpdateRoom(ctx context.Context, param UpdateRoom) (Room, error) {
	rm, err := rs.Querier.UpdateRoom(ctx, db.UpdateRoomParams{
		Status: db.RoomStatus(param.Status),
		ID: pgtype.UUID{
			Bytes: param.ID,
			Valid: true,
		},
	})
	if err != nil {
		rs.logger.Error("failed to update room", err)
		return Room{}, err
	}
	return Room{
		ID:         rm.ID.Bytes,
		RoomNumber: rm.RoomNumber,
		HotelID:    rm.HotelID.Bytes,
		CreatedAt:  rm.CreatedAt.Time,
		UpdatedAt:  rm.UpdatedAt.Time,
	}, nil
}
func (rs *roomService) GetRoomReservations(ctx context.Context, roomID string) ([]db.Reservation, error) {
	//TODO:add filter to get rooms
	rvns, err := rs.Querier.GetRoomReservations(ctx,
		pgtype.UUID{
			Bytes: uuid.MustParse(roomID),
			Valid: true,
		})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			rs.logger.Info("rooms not found", err)
			return nil, &apperror.AppError{ErrorCode: http.StatusNotFound, RootError: apperror.ErrRecordNotFound}
		}
		rs.logger.Error("failed to get rooms", err)
		return nil, &apperror.AppError{ErrorCode: http.StatusInternalServerError, RootError: apperror.ErrUnableToGet}
	}

	return rvns, nil
}

// search rooms based on price,location,date
// and other criterias

func (rs *roomService) SearchRoom(ctx context.Context, searchParam SearchParam) ([]db.Room, error) {
	if err := searchParam.Validate(); err != nil {
		rs.logger.Info("invalid search param", err)
		return nil, &apperror.AppError{ErrorCode: http.StatusBadRequest, RootError: err}
	}
	rooms, err := rs.Querier.SearchRoom(ctx, db.SearchRoomParams{
		Price:         searchParam.Price,
		StGeogpoint:   searchParam.Location.Latitude,
		StGeogpoint_2: searchParam.Location.Longitude,
		FromTime: pgtype.Timestamptz{
			Time:  searchParam.FromTime,
			Valid: true,
		},
		FromTime_2: pgtype.Timestamptz{
			Time:  searchParam.ToTime,
			Valid: true,
		},
		RoomType: db.Roomtype(searchParam.RoomType),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			rs.logger.Info("rooms not found", err)
			return nil, &apperror.AppError{ErrorCode: http.StatusNotFound, RootError: apperror.ErrRecordNotFound}
		}
		rs.logger.Error("failed to get rooms", err)
		return nil, &apperror.AppError{ErrorCode: http.StatusInternalServerError, RootError: apperror.ErrUnableToGet}
	}

	return rooms, nil
}
func (rs *roomService) AddRoom(ctx context.Context, param CreateRoomParam) (CreateRoomResponse, error) {
	conn, err := rs.Pool.Acquire(ctx)
	if err != nil {
		rs.logger.Error("failed to acquire connection for transaction", err)
		return CreateRoomResponse{}, &apperror.AppError{
			ErrorCode: http.StatusInternalServerError,
			RootError: errors.New("failed to acquire db connection"),
		}
	}
	defer conn.Release()
	queries := db.New(conn)
	tx, err := conn.Begin(ctx)
	if err != nil {
		rs.logger.Error("failed to create tx instance", err)
		return CreateRoomResponse{}, &apperror.AppError{
			ErrorCode: http.StatusInternalServerError,
			RootError: errors.New("failed to add room"),
		}
	}

	qtx := queries.WithTx(tx)
	roomType, err := qtx.AddRoomType(ctx, db.AddRoomTypeParams{
		RoomType:     param.RoomTypeParam.RoomType,
		Price:        param.RoomTypeParam.Price,
		Description:  param.RoomTypeParam.Description,
		MaxAccupancy: param.RoomTypeParam.MaxAccupancy,
	})
	if err != nil {
		rs.logger.Error("failed to add room", err)
		return CreateRoomResponse{}, &apperror.AppError{
			ErrorCode: http.StatusInternalServerError,
			RootError: apperror.ErrUnableToCreate}
	}
	room, err := qtx.AddRoom(ctx, db.AddRoomParams{
		RoomNumber: param.RoomParam.RoomNumber,
		HotelID:    param.RoomParam.HotelID,
		RoomTypeID: roomType.ID,
		Floor:      param.RoomParam.Floor,
	})
	if err != nil {
		rs.logger.Error("failed to add room", err)
		return CreateRoomResponse{}, &apperror.AppError{
			ErrorCode: http.StatusInternalServerError,
			RootError: apperror.ErrUnableToCreate,
		}
	}
	return CreateRoomResponse{
		Room:     room,
		RoomType: roomType,
	}, nil
}
