package room

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	apperror "reservation/internal/app_error"
	"reservation/internal/storage/db"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/paymentintent"

	"golang.org/x/exp/slog"
)

type RoomService interface {
	ReserveRoom(ctx context.Context, param ReserveRoom) (string, error)
	UpdateRoom(ctx context.Context, param UpdateRoom) (Room, error)
	WebhookAction(ctx context.Context, event stripe.Event) error
	GetRoomReservations(ctx context.Context, roomID string) ([]db.Reservation, error)
}

type ReservationStatus string

const (
	StatusPending    ReservationStatus = "PENDING"
	StatusSuccessful ReservationStatus = "SUCCESSFUL"
	StatusFailed     ReservationStatus = "FAILED"
)

type roomService struct {
	db.Querier
	logger          *slog.Logger
	stripeSecretKey string
}

func NewRoomService(q db.Querier, logger *slog.Logger, key string) RoomService {
	return &roomService{
		Querier:         q,
		stripeSecretKey: key,
		logger:          logger,
	}
}

func (rs *roomService) ReserveRoom(ctx context.Context, param ReserveRoom) (string, error) {
	if err := param.Validate(); err != nil {

		return "", &apperror.AppError{
			ErrorCode: http.StatusBadRequest,
			RootError: err,
		}
	}

	rvn, err := rs.CreateReservation(ctx, db.CreateReservationParams{
		RoomID:   pgtype.UUID{Bytes: param.RoomID, Valid: true},
		UserID:   pgtype.UUID{Bytes: param.UserID, Valid: true},
		Status:   db.ReservationStatus(StatusPending),
		FromTime: pgtype.Timestamptz{Time: param.FromTime, Valid: true},
		ToTime:   pgtype.Timestamptz{Time: param.ToTime, Valid: true},
	})
	if err != nil {
		rs.logger.Error("failed to create reservation", err)
		return "", &apperror.AppError{
			ErrorCode: http.StatusInternalServerError,
			RootError: errors.New("failed to make reservation"),
		}
	}
	secretKey, err := rs.createPaymentIntent(ctx, rvn.ID.String(), param.RoomID.String())
	if err != nil {
		return "", err
	}

	return secretKey, nil

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
func (rs *roomService) createPaymentIntent(ctx context.Context, rvnID, roomID string) (string, error) {
	room, err := rs.Querier.GetRoom(ctx, pgtype.UUID{Bytes: uuid.MustParse(roomID), Valid: true})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			rs.logger.Info("room not found", err)
			return "", &apperror.AppError{ErrorCode: http.StatusNotFound, RootError: ErrRecordNotFound}
		}
		rs.logger.Error("failed to get room", err)
		return "", &apperror.AppError{ErrorCode: http.StatusInternalServerError, RootError: ErrUnableToGet}
	}

	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(int64(room.Price)),
		Currency: stripe.String(string(stripe.CurrencyUSD)),
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
		Metadata: map[string]string{
			"room_id": roomID,
			"id":      rvnID,
		},
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		rs.logger.Error("failed to create stripe payment intent")
		return "", &apperror.AppError{
			ErrorCode: http.StatusInternalServerError,
			RootError: err,
		}
	}

	return pi.ClientSecret, nil
}

func (rs *roomService) WebhookAction(ctx context.Context, event stripe.Event) error {
	switch event.Type {
	case "payment_intent.succeeded":
		var paymentIntent stripe.PaymentIntent
		err := json.Unmarshal(event.Data.Raw, &paymentIntent)
		if err != nil {
			rs.logger.Error("Error parsing webhook JSON", err)
			return &apperror.AppError{
				ErrorCode: http.StatusInternalServerError,
				RootError: errors.New("failed to unmarshal data"),
			}
		}
		rvn, err := rs.UpdateReservation(ctx, db.UpdateReservationParams{
			Status: db.ReservationStatus(StatusSuccessful),
			ID: pgtype.UUID{
				Bytes: uuid.MustParse(paymentIntent.Metadata["id"]),
				Valid: true,
			},
		})
		if err != nil {
			rs.logger.Error("failed to update reservation", err, rvn, paymentIntent)
			return err
		}
	case "payment_method.attached":
		var paymentMethod stripe.PaymentMethod
		err := json.Unmarshal(event.Data.Raw, &paymentMethod)
		if err != nil {
			rs.logger.Error("Error parsing webhook JSON", err, paymentMethod)
			return err
		}

	default:
		rs.logger.Info("unhandled envet type", event.Type)
		return nil
	}

	return nil
	// TODO: change status of reservation to FAILED
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
			return nil, &apperror.AppError{ErrorCode: http.StatusNotFound, RootError: ErrRecordNotFound}
		}
		rs.logger.Error("failed to get rooms", err)
		return nil, &apperror.AppError{ErrorCode: http.StatusInternalServerError, RootError: ErrUnableToGet}
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
			return nil, &apperror.AppError{ErrorCode: http.StatusNotFound, RootError: ErrRecordNotFound}
		}
		rs.logger.Error("failed to get rooms", err)
		return nil, &apperror.AppError{ErrorCode: http.StatusInternalServerError, RootError: ErrUnableToGet}
	}

	return rooms, nil
}
