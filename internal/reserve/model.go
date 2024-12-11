package reserve

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
)

type ReserveRequest struct {
	RoomID uuid.UUID `json:"room_id"`
	UserID uuid.UUID `json:"user_id"`
}
type ReserveResponse struct {
	ID         uuid.UUID
	RoomNumber string
	UserID     uuid.UUID
	HotelID    uuid.UUID
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
type CheckoutResponse struct {
	SessionID  uuid.UUID `json:"session_id"`
	PaymentURL string    `json:"payment_url"`
	CancelURL  string    `json:"cancel_url"`
}

func (rr ReserveRequest) Validate() error {
	return validation.ValidateStruct(
		&rr,
		validation.Field(&rr.RoomID, validation.Required.Error("room id is required")),
		validation.Field(&rr.UserID, validation.Required.Error("user id is required")),
	)
}

type CheckoutRequest struct {
	ProductID   string
	CallbackURL string `json:"callback_url"`
}

type CallBackRequest struct {
}
