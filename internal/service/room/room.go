package room

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"reservation/internal/storage/db"

	"github.com/google/uuid"
	"golang.org/x/exp/slog"
)

type Reserver interface {
	ReserveRoom(ctx context.Context, param ReserveRoom) (string, error)
	UpdateRoom(ctx context.Context, param UpdateRoom) (Room, error)
}

type room struct {
	db.Querier
	url    string
	logger slog.Logger
}

func Init(q db.Querier, url string) Reserver {
	return &room{
		Querier: q,
		url:     url,
	}
}

func (r *room) ReserveRoom(ctx context.Context, param ReserveRoom) (string, error) {
	if err := param.Validate(); err != nil {
		return "", ErrInvalidInput
	}
	_, err := r.Querier.UpdateRoom(ctx, db.UpdateRoomParams{
		UserID: uuid.NullUUID{
			UUID:  param.UserID,
			Valid: true,
		},

		ID: param.RoomID,
	})
	if err != nil {
		return "", ErrReservationFailed
	}
	req := CheckoutRequest{
		ProductID:   param.RoomID.String(),
		CallbackURL: "http://localhost:9090/callback", //TODO: url should be read from config
	}
	ssn, err := r.createCheckoutSession(ctx, req)
	if err != nil {
		return "", ErrCheckoutSessionFailed
	}
	return ssn.PaymentURL, nil

}
func (r *room) createCheckoutSession(ctx context.Context, req CheckoutRequest) (CheckoutResponse, error) {
	bbyte, err := json.Marshal(req)
	if err != nil {
		return CheckoutResponse{}, err
	}

	request, err := http.NewRequest(http.MethodPost, r.url, bytes.NewBuffer(bbyte))
	if err != nil {
		return CheckoutResponse{}, err
	}
	client := http.Client{}
	res, err := client.Do(request)
	if err != nil {
		return CheckoutResponse{}, err
	}
	var session CheckoutResponse
	if err := json.NewDecoder(res.Body).Decode(&session); err != nil {
		return CheckoutResponse{}, err
	}
	return session, nil
}
func (r *room) UpdateRoom(ctx context.Context, param UpdateRoom) (Room, error) {
	rm, err := r.Querier.UpdateRoom(ctx, db.UpdateRoomParams{
		Status: db.RoomStatus(param.Status),
		UserID: uuid.NullUUID{
			UUID:  param.UserID,
			Valid: true,
		},
		ID: param.ID,
	})
	if err != nil {
		r.logger.Error("failed to update room", err)
		return Room{}, err
	}
	return Room{
		ID:         rm.ID,
		RoomNumber: rm.RoomNumber,
		UserID:     rm.UserID.UUID,
		HotelID:    rm.HotelID,
		CreatedAt:  rm.CreatedAt,
		UpdatedAt:  rm.UpdatedAt,
	}, nil
}
