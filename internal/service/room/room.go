package room

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"reservation/internal/storage/db"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stripe/stripe-go/paymentintent"
	"github.com/stripe/stripe-go/v78"

	"golang.org/x/exp/slog"
)

type RoomService interface {
	ReserveRoom(ctx context.Context, param ReserveRoom) (string, error)
	UpdateRoom(ctx context.Context, param UpdateRoom) (Room, error)
	WebhookAction(ctx context.Context, event stripe.Event)
}

type ReservationStatus string

const (
	StatusPending    ReservationStatus = "PENDING"
	StatusSuccessful ReservationStatus = "SUCCESSFUL"
	StatusFailed     ReservationStatus = "FAILED"
)

type roomService struct {
	db.Querier
	url             string
	logger          slog.Logger
	stripeSecretKey string
}

func NewRoomService(q db.Querier, url string) RoomService {
	return &roomService{
		Querier: q,
		url:     url,
	}
}

func (rs *roomService) ReserveRoom(ctx context.Context, param ReserveRoom) (string, error) {
	if err := param.Validate(); err != nil {
		return "", ErrInvalidInput
	}

	_, err := rs.CreateReservation(ctx, db.CreateReservationParams{
		RoomID:   pgtype.UUID{Bytes: param.RoomID, Valid: true},
		UserID:   pgtype.UUID{Bytes: param.UserID, Valid: true},
		Status:   db.ReservationStatus(StatusPending),
		FromTime: pgtype.Timestamptz{Time: param.FromTime, Valid: true},
		ToTime:   pgtype.Timestamptz{Time: param.ToTime, Valid: true},
	})
	if err != nil {
		return "", errors.New("failed to make reservation")
	}

	req := CheckoutRequest{
		ProductID:   param.RoomID.String(),
		CallbackURL: "http://localhost:9090/callback", //TODO: url should be read from config
	}
	ssn, err := rs.createCheckoutSession(ctx, req)
	if err != nil {
		return "", ErrCheckoutSessionFailed
	}
	return ssn.PaymentURL, nil

}
func (rs *roomService) createCheckoutSession(ctx context.Context, req CheckoutRequest) (CheckoutResponse, error) {
	bbyte, err := json.Marshal(req)
	if err != nil {
		return CheckoutResponse{}, err
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, rs.url, bytes.NewBuffer(bbyte))
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
func (rs *roomService) UpdateRoom(ctx context.Context, param UpdateRoom) (Room, error) {
	rm, err := rs.Querier.UpdateRoom(ctx, db.UpdateRoomParams{
		Status: db.RoomStatus(param.Status),
		UserID: pgtype.UUID{
			Bytes: param.UserID,
			Valid: true,
		},

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
		UserID:     rm.UserID.Bytes,
		HotelID:    rm.HotelID.Bytes,
		CreatedAt:  rm.CreatedAt.Time,
		UpdatedAt:  rm.UpdatedAt.Time,
	}, nil
}
func (rs *roomService) createPaymentIntent(ctx, roomID, userID string) (string, error) {

	stripe.Key = "sk_test_4eC39HqLyjWDarjtT1zdp7dc" // change this key to my own

	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(2000),
		Currency: stripe.String(string(stripe.CurrencyUSD)),
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
	}
	result, err := paymentintent.New(params)
	return "", nil
}

func (rs *roomService) WebhookAction(ctx context.Context, event stripe.Event) {
	switch event.Type {
	case "payment_intent.succeeded":
		var paymentIntent stripe.PaymentIntent
		err := json.Unmarshal(event.Data.Raw, &paymentIntent)
		if err != nil {
			rs.logger.Error("Error parsing webhook JSON", err)
			return
		}
		// change status of reservation to SUCCESSFUL
	case "payment_method.attached":
		var paymentMethod stripe.PaymentMethod
		err := json.Unmarshal(event.Data.Raw, &paymentMethod)
		if err != nil {
			rs.logger.Error("Error parsing webhook JSON", err)
			return
		}

	default:
		rs.logger.Info("unhandled envet type", event.Type)
	}

	// TODO: change status of reservation to FAILED
}
