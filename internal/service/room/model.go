package room

import (
	"reservation/internal/storage/db"
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
	RoomNumber int32
	UserID     uuid.UUID
	HotelID    uuid.UUID
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type ReserveRoom struct {
	RoomID      uuid.UUID `json:"room_id"`
	UserID      uuid.UUID `json:"user_id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	PhoneNumber string    `json:"phone_number"`
	Email       string    `json:"email"`
	FromTime    time.Time `json:"from_time"`
	ToTime      time.Time `json:"to_time"`
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
		validation.Field(&rr.FirstName, validation.Required.Error("fist_name is required")),
		validation.Field(&rr.LastName, validation.Required.Error("last_name is required")),
		validation.Field(&rr.PhoneNumber, validation.Required.Error("phone_number is required")),
		validation.Field(&rr.Email, validation.Required.Error("email is required")),
		validation.Field(&rr.FromTime, validation.Required.Error("From time is required"),
			validation.Min(time.Now()).Error("from time can not be past time")),
	)
}

type CheckoutRequest struct {
	ProductID   string `json:"product_id"`
	CallbackURL string `json:"callback_url"`
}

type CallBackRequest struct {
}
type SearchParam struct {
	FromTime time.Time `json:"From_time"`
	ToTime   time.Time `json:"to_time"`
	Location struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"location"`
	Price    float64 `json:"price"`
	RoomType string  `json:"room_type"`
}

func (sp SearchParam) Validate() error {
	return validation.ValidateStruct(&sp,
		validation.Field(&sp.FromTime, validation.Required.Error("from_time is required")),
		validation.Field(&sp.ToTime, validation.Required.Error("to_time is required")),
	)
}

type CreateRoomParam struct {
	RoomTypeParam db.AddRoomTypeParams `json:"room_type_param"`
	RoomParam     db.AddRoomParams     `json:"room_param"`
}
type CreateRoomResponse struct {
	Room     db.Room     `json:"room"`
	RoomType db.RoomType `json:"room_type"`
}

type GetHotelRoomsResponse struct {
	Room     db.Room     `json:"room"`
	RoomType db.RoomType `json:"room_type"`
	Count    int64       `json:"count"`
}
