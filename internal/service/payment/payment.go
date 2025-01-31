package payment

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
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
	HandleWebHook(c *gin.Context)
}
type PaymentProviderConfig struct {
	BaseURL      string
	ReturnURL    string
	CancelURL    string
	ClientID     string
	ClientSecret string
}
type paymentService struct {
	logger       *slog.Logger
	paypalConfig PaymentProviderConfig
	db.Querier
}

func NewPaymentService(logger *slog.Logger, q db.Querier, config PaymentProviderConfig) PaymentProcessor {
	return &paymentService{
		logger:       logger,
		Querier:      q,
		paypalConfig: config,
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
		paymentURL, err = p.createPaypalPaymentIntent(context.Background(), rvn.ID.String(), rvn.RoomID.String())
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

func (p *paymentService) createPaypalPaymentIntent(ctx context.Context, resID, roomID string) (string, error) {

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

	accessToken, err := p.getPaypalAccessToken()
	if err != nil {
		log.Fatalf("Failed to get access token: %v", err)
	}

	// Step 2: Create a Payment
	paymentID, err := p.createPaypalPayment(accessToken, customData{
		ReservationID: resID,
		RoomID:        room.ID.String(),
		Price:         fmt.Sprintf("%0.2f", room.Price),
	})

	if err != nil {
		log.Fatalf("Failed to create payment: %v", err)
	}

	fmt.Println("Payment created with ID:", paymentID)

	return "", nil
}

func (p *paymentService) getPaypalAccessToken() (string, error) {
	url := p.paypalConfig.BaseURL + "/v1/oauth2/token"
	req, err := http.NewRequest("POST", url, bytes.NewBufferString("grant_type=client_credentials"))
	if err != nil {
		return "", err
	}

	req.SetBasicAuth(p.paypalConfig.ClientID, p.paypalConfig.ClientSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var tokenResponse AccessTokenResponse
	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		return "", err
	}

	return tokenResponse.AccessToken, nil
}

func (p *paymentService) createPaypalPayment(accessToken string, customData customData) (string, error) {
	url := p.paypalConfig.BaseURL + "/v2/checkout/orders"

	// Define metadata including customer_id
	metadata := map[string]string{
		"reservation_id": customData.ReservationID,
		"room_id":        customData.RoomID,
	}
	metadataJSON, _ := json.Marshal(metadata)

	orderRequest := map[string]interface{}{
		"intent": "CAPTURE",
		"purchase_units": []map[string]interface{}{
			{
				"amount": map[string]string{
					"currency_code": "USD",
					"value":         customData.Price,
				},
				"description": "Payment for room reservation",
				"custom_id":   string(metadataJSON), // Include customer_id and other metadata
			},
		},
		"application_context": map[string]string{
			"return_url": p.paypalConfig.ReturnURL,
			"cancel_url": p.paypalConfig.CancelURL,
		},
	}

	jsonData, err := json.Marshal(orderRequest)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var orderResponse CreateOrderResponse
	if err := json.Unmarshal(body, &orderResponse); err != nil {
		return "", err
	}

	return orderResponse.ID, nil
}

func (p *paymentService) HandleWebHook(c *gin.Context) {
	switch {
	case c.Request.Header.Get("Stripe-Signature") != "":
		buf := make([]byte, 0)
		sigHeader := c.Request.Header.Get("Stripe-Signature")
		_, err := c.Request.Body.Read(buf)
		if err != nil {
			p.logger.Info("fialed to read stripe webhook payload", err)
			return
		}

		event, err := webhook.ConstructEvent(buf, sigHeader, "secret")
		if err != nil {
			p.logger.Info("fialed to bind stripe webhook body", err)
			return
		}
		p.HandleStripeWebHook(context.Background(), event)

	case c.Request.Header.Get("X-Razorpay-Signature") != "":
		// handle razorpay webhook

	case c.Request.Header.Get("PayPal-Transmission-Sig") != "":
		//handle paypal webhook

	}

}
func (p *paymentService) HandleStripeWebHook(ctx context.Context, event stripe.Event) {
	switch event.Type {
	case "payment_intent.succeeded":
		var paymentIntent stripe.PaymentIntent
		err := json.Unmarshal(event.Data.Raw, &paymentIntent)
		if err != nil {
			p.logger.Error("Error parsing webhook JSON", err)
			return
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
		}
	case "payment_method.attached":
		var paymentMethod stripe.PaymentMethod
		err := json.Unmarshal(event.Data.Raw, &paymentMethod)
		if err != nil {
			p.logger.Error("Error parsing webhook JSON", err, paymentMethod)
		}

	default:
		p.logger.Info("unhandled envet type", event.Type)
	}

	// TODO: change status of reservation to FAILED
}
