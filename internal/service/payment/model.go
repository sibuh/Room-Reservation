package payment

import "encoding/json"

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

type CreateOrderResponse struct {
	ID     string `json:"id"`
	Links  []Link `json:"links"`
	Status string `json:"status"`
}

type Link struct {
	Href   string `json:"href"`
	Rel    string `json:"rel"`
	Method string `json:"method"`
}
type customData struct {
	ReservationID string `json:"reservation_id"`
	RoomID        string `json:"room_id"`
	Price         string `json:"price"`
}
type CaptureOrderResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

type PaypalWebhookPayload struct {
	EventType string `json:"event_type"`
	Resource  struct {
		ID       string `json:"id"`
		CustomID string `json:"custom_id"`
		Status   string `json:"status"`
	} `json:"resource"`
}

type VerifyWebhookRequest struct {
	TransmissionID   string          `json:"transmission_id"`
	TransmissionTime string          `json:"transmission_time"`
	TransmissionSig  string          `json:"transmission_sig"`
	CertURL          string          `json:"cert_url"`
	WebhookID        string          `json:"webhook_id"`
	WebhookEvent     json.RawMessage `json:"webhook_event"`
}
type VerifyWebhookResponse struct {
	VerificationStatus string `json:"verification_status"`
}
