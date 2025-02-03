package payment

import (
	"bytes"
	"context"
	"crypto/sha256"
	"crypto/x509"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
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
	CapturePaypalPayment(ctx context.Context, orderID string) error
	HandleWebHook(c *gin.Context)
}
type PaymentProviderConfig struct {
	BaseURL      string
	ReturnURL    string
	CancelURL    string
	ClientID     string
	ClientSecret string
	WebHookID    string
	StripeSecret string
}
type paymentService struct {
	logger       *slog.Logger
	paypalConfig PaymentProviderConfig
	StripeConfig PaymentProviderConfig
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
			"id":      rvnID,
			"room_id": roomID,
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
		p.logger.Error("Failed to get paypal access token: %v", err)
		return "", &apperror.AppError{
			ErrorCode: http.StatusInternalServerError,
			RootError: errors.New("failed to get paypal access token"),
		}

	}

	// Step 2: Create a Payment
	paymentID, err := p.createPaypalPayment(accessToken, customData{
		ReservationID: resID,
		RoomID:        room.ID.String(),
		Price:         fmt.Sprintf("%0.2f", room.Price),
	})

	if err != nil {
		p.logger.Error("failed to create paypal payment", err)
		return "", &apperror.AppError{
			ErrorCode: http.StatusInternalServerError,
			RootError: errors.New("failed to create paypal payment intent"),
		}

	}

	return paymentID, nil
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

	body, err := io.ReadAll(resp.Body)
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
		"id":      customData.ReservationID,
		"room_id": customData.RoomID,
	}
	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		p.logger.Error("failed to marshal the metadata to be added in paypal createPayment", err)
		return "", errors.New("marshal error")
	}

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
		p.logger.Error("failed to marshal paypal orderRequest", err)
		return "", &apperror.AppError{
			ErrorCode: http.StatusInternalServerError,
			RootError: errors.New("marshal error"),
		}
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var orderResponse CreateOrderResponse
	if err := json.Unmarshal(body, &orderResponse); err != nil {
		return "", err
	}
	for _, link := range orderResponse.Links {
		if link.Rel == "approve" {
			return link.Href, nil
		}
	}

	return "", nil
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

		event, err := webhook.ConstructEvent(buf, sigHeader, p.StripeConfig.StripeSecret)
		if err != nil {
			p.logger.Info("fialed to bind stripe webhook body", err)
			return
		}
		p.HandleStripeWebHook(context.Background(), event)

	case c.Request.Header.Get("X-Razorpay-Signature") != "":
		// handle razorpay webhook

	case c.Request.Header.Get("PayPal-Transmission-Sig") != "":
		//handle paypal webhook

		// Step 1: Extract headers and payload
		transmissionID := c.GetHeader("Paypal-Transmission-Id")
		transmissionTime := c.GetHeader("Paypal-Transmission-Time")
		transmissionSig := c.GetHeader("Paypal-Transmission-Sig")
		certURL := c.GetHeader("Paypal-Cert-Url")

		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			p.logger.Error("failed to read paypal webhook payload", err)
			return
		}

		// Step 2: Verify the webhook
		isValid, err := p.verifyWebhookAPI(transmissionID, transmissionTime, transmissionSig, certURL, body)
		if err != nil {
			p.logger.Error("paypal webhook request verification failed", err)
			return
		}

		if !isValid {
			p.logger.Error("paypal webhook request is not valid", errors.New("invalid paypal webhook request"))
			return
		}

		// Step 3: Process the webhook payload
		var event PaypalWebhookPayload
		if err := json.Unmarshal(body, &event); err != nil {
			p.logger.Error("Failed to unmarshal paypal webhook payload", err)
			return
		}
		p.HandlePaypalWebHook(context.Background(), event)
	}

}
func (p *paymentService) verifyWebhookAPI(transmissionID, transmissionTime, transmissionSig, certURL string, body []byte) (bool, error) {
	// Step 1: Fetch the PayPal public certificate
	cert, err := fetchPayPalCertificate(certURL)
	if err != nil {
		return false, fmt.Errorf("failed to fetch PayPal certificate: %v", err)
	}

	// Step 2: Decode the signature
	signature, err := base64.StdEncoding.DecodeString(transmissionSig)
	if err != nil {
		return false, fmt.Errorf("failed to decode signature: %v", err)
	}

	// Step 3: Create the signed data string
	signedData := fmt.Sprintf("%s|%s|%s|%s", transmissionID, transmissionTime, p.paypalConfig.WebHookID, string(body))

	// Step 4: Hash the signed data
	hashed := sha256.Sum256([]byte(signedData))

	// Step 5: Verify the signature
	err = cert.CheckSignature(x509.SHA256WithRSA, hashed[:], signature)
	if err != nil {
		return false, fmt.Errorf("signature verification failed: %v", err)
	}

	return true, nil

}

func fetchPayPalCertificate(certURL string) (*x509.Certificate, error) {
	resp, err := http.Get(certURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch certificate: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read certificate body: %v", err)
	}

	// Decode the PEM-encoded certificate
	block, _ := pem.Decode(body)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate: %v", err)
	}

	return cert, nil
}

func (p *paymentService) CapturePaypalPayment(ctx context.Context, orderID string) error {
	token, err := p.getPaypalAccessToken()
	if err != nil {
		p.logger.Error("failed to get paypal access token", err)
		return &apperror.AppError{
			ErrorCode: http.StatusInternalServerError,
			RootError: apperror.ErrUnableToCreate,
		}
	}
	_, err = p.captureOrderPayment(token, orderID)
	if err != nil {
		p.logger.Error("failed to capture paypal payment", err)
		return &apperror.AppError{
			ErrorCode: http.StatusInternalServerError,
			RootError: apperror.ErrCapturingPayment,
		}
	}
	return nil
}

func (p *paymentService) captureOrderPayment(accessToken, orderID string) (CaptureOrderResponse, error) {
	url := p.paypalConfig.BaseURL + "/v2/checkout/orders/" + orderID + "/capture"

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return CaptureOrderResponse{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return CaptureOrderResponse{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return CaptureOrderResponse{}, err
	}

	var captureResponse CaptureOrderResponse
	if err := json.Unmarshal(body, &captureResponse); err != nil {
		return CaptureOrderResponse{}, err
	}

	return captureResponse, nil
}

func (p *paymentService) HandleStripeWebHook(ctx context.Context, event stripe.Event) {
	switch event.Type {
	case "payment_intent.succeeded":
		var paymentIntent stripe.PaymentIntent
		err := json.Unmarshal(event.Data.Raw, &paymentIntent)
		if err != nil {
			p.logger.Error("Error parsing stripe webhook payload", err)
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

func (p *paymentService) HandlePaypalWebHook(ctx context.Context, event PaypalWebhookPayload) {
	// Extract metadata from the webhook payload
	customID := event.Resource.CustomID
	// status := event.Resource.Status

	// Parse custom_id (metadata)
	var metadata map[string]string
	if err := json.Unmarshal([]byte(customID), &metadata); err != nil {
		p.logger.Error("Failed to parse custom_id:", err)
		return
	}
	// update reservation status
	rvn, err := p.UpdateReservation(ctx, db.UpdateReservationParams{
		Status: db.ReservationStatus(room.StatusSuccessful),
		ID: pgtype.UUID{
			Bytes: uuid.MustParse(metadata["id"]),
			Valid: true,
		},
	})
	if err != nil {
		p.logger.Error("failed to update reservation", err, rvn, metadata["id"])
	}

}
