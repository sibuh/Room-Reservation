package payment

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"reservation/internal/apperror"
	"reservation/internal/service/room"
	"reservation/internal/storage/db"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/paymentintent"
	"github.com/stripe/stripe-go/v78/webhook"
	"golang.org/x/exp/slog"
)

type PaymentProcessor interface {
	ProcessPayment(ctx context.Context, agent string, rvn db.Reservation) (string, error)
	HandleWebHook(c *gin.Context) error
}
type paymentService struct {
	logger *slog.Logger
	db.Querier
}

func NewPaymentService(logger *slog.Logger, q db.Querier) PaymentProcessor {
	return &paymentService{
		logger:  logger,
		Querier: q,
	}
}

func (p *paymentService) ProcessPayment(ctx context.Context, agent string, rvn db.Reservation) (string, error) {
	var paymentURL string
	var err error
	switch agent {
	case "stripe":
		//process payment with stripe
		paymentURL, err = p.createStripePaymentIntent(context.Background(), rvn.ID.String(), rvn.RoomID.String())
	case "paypal":
		//process payment with paypal
	case "telebirr":
		//process payment with telebirr
	case "chapa":
		//process payment with chapa
	}
	return paymentURL, err
}

func (p *paymentService) createStripePaymentIntent(ctx context.Context, rvnID, roomID string) (string, error) {
	room, err := p.Querier.GetRoom(ctx, pgtype.UUID{Bytes: uuid.MustParse(roomID), Valid: true})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			p.logger.Info("room not found", err)
			return "", &apperror.AppError{ErrorCode: http.StatusNotFound, RootError: apperror.ErrRecordNotFound}
		}
		p.logger.Error("failed to get room", err)
		return "", &apperror.AppError{
			ErrorCode: http.StatusInternalServerError,
			RootError: apperror.ErrUnableToGet}
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
		p.logger.Error("failed to create stripe payment intent")
		return "", &apperror.AppError{
			ErrorCode: http.StatusInternalServerError,
			RootError: err,
		}
	}

	return pi.ClientSecret, nil
}

func (p *paymentService) HandleWebHook(c *gin.Context) error {
	switch {
	case c.Request.Header.Get("Stripe-Signature") != "":
		payload := make([]byte, 0)
		sigHeader := c.Request.Header.Get("Stripe-Signature")
		_, err := c.Request.Body.Read(payload)
		if err != nil {
			p.logger.Info("fialed to read stripe webhook payload", err)
			_ = c.Error(&apperror.AppError{
				ErrorCode: http.StatusBadRequest,
				RootError: apperror.ErrInvalidInput,
			})
		}
		event, err := webhook.ConstructEvent(payload, sigHeader, "secret")
		if err != nil {
			p.logger.Info("fialed to bind stripe webhook body", err)
			_ = c.Error(&apperror.AppError{
				ErrorCode: http.StatusBadRequest,
				RootError: apperror.ErrBindingRequestBody,
			})
		}
		return p.HandleStripeWebHook(context.Background(), event)

	case c.Request.Header.Get("X-Razorpay-Signature") != "":
		// handle razorpay webhook

	case c.Request.Header.Get("PayPal-Transmission-Sig") != "":
		//handle paypal webhook

	}

	return nil

}
func (p *paymentService) HandleStripeWebHook(ctx context.Context, event stripe.Event) error {
	switch event.Type {
	case "payment_intent.succeeded":
		var paymentIntent stripe.PaymentIntent
		err := json.Unmarshal(event.Data.Raw, &paymentIntent)
		if err != nil {
			p.logger.Error("Error parsing webhook JSON", err)
			return &apperror.AppError{
				ErrorCode: http.StatusInternalServerError,
				RootError: errors.New("failed to unmarshal data"),
			}
		}
		rvn, err := p.UpdateReservation(ctx, db.UpdateReservationParams{
			Status: db.ReservationStatus(room.StatusSuccessful),
			ID: pgtype.UUID{
				Bytes: uuid.MustParse(paymentIntent.Metadata["id"]),
				Valid: true,
			},
		})
		if err != nil {
			p.logger.Error("failed to update reservation", err, rvn, paymentIntent)
			return err
		}
	case "payment_method.attached":
		var paymentMethod stripe.PaymentMethod
		err := json.Unmarshal(event.Data.Raw, &paymentMethod)
		if err != nil {
			p.logger.Error("Error parsing webhook JSON", err, paymentMethod)
			return err
		}

	default:
		p.logger.Info("unhandled envet type", event.Type)
		return nil
	}

	return nil
	// TODO: change status of reservation to FAILED
}
