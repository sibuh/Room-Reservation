package room

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"reservation/internal/storage/db"

	"github.com/google/uuid"
)

type reserve struct {
	db.Querier
	url string
}

func Init(q db.Querier) reserve {
	return reserve{
		Querier: q,
	}
}

func (r *reserve) ReserveRoom(ctx context.Context, param ReserveRequest) (CheckoutResponse, error) {
	if err := param.Validate(); err != nil {
		return CheckoutResponse{}, ErrInvalidInput
	}
	_, err := r.Querier.HoldRoom(ctx, db.HoldRoomParams{
		UserID: uuid.NullUUID{
			UUID:  param.UserID,
			Valid: true,
		},

		ID: param.RoomID,
	})
	if err != nil {
		return CheckoutResponse{}, ErrReservationFailed
	}
	req := CheckoutRequest{
		ProductID:   param.RoomID.String(),
		CallbackURL: "http://localhost:9090/callback", //TODO: url should be read from config
	}
	ssn, err := r.createCheckoutSession(ctx, req)
	if err != nil {
		return CheckoutResponse{}, ErrCheckoutSessionFailed
	}
	return ssn, nil

}
func (r *reserve) createCheckoutSession(ctx context.Context, req CheckoutRequest) (CheckoutResponse, error) {
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
func (r *reserve) UpdateRoomStatus(ctx context.Context, cbr CallBackRequest) (interface{}, error) {
	return nil, nil
}
