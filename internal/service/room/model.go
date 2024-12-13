package room

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
)

type RoomStatus string

const (
	Free     RoomStatus = "FREE"
	Reserved RoomStatus = "RESERVED"
	Held     RoomStatus = "HELD"
)

type Room struct {
	ID         uuid.UUID
	RoomNumber string
	UserID     uuid.UUID
	HotelID    uuid.UUID
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type ReserveRoom struct {
	RoomID uuid.UUID `json:"room_id"`
	UserID uuid.UUID `json:"user_id"`
}
type UpdateRoom struct {
	ID     uuid.UUID
	Status RoomStatus
	UserID uuid.UUID
}
type CheckoutResponse struct {
	SessionID  uuid.UUID `json:"session_id"`
	PaymentURL string    `json:"payment_url"`
	CancelURL  string    `json:"cancel_url"`
}

func (rr ReserveRoom) Validate() error {
	return validation.ValidateStruct(
		&rr,
		validation.Field(&rr.RoomID, validation.Required.Error("room id is required")),
		validation.Field(&rr.UserID, validation.Required.Error("user id is required")),
	)
}

type CheckoutRequest struct {
	ProductID   string `json:"product_id"`
	CallbackURL string `json:"callback_url"`
}

type CallBackRequest struct {
}
