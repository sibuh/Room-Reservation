package reserve

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

type ReserveRequest struct {
	HotelID uuid.UUID `json:"hotel_id"`
	RoomID  uuid.UUID `json:"room_id"`
	UserID  uuid.UUID `json:"user_id"`
}
type Session struct {
	SessionID  uuid.UUID `json:"session_id"`
	PaymentURL string    `json:"payment_url"`
	CancelURL  string    `json:"cancel_url"`
}

var ErrReservationFailed = errors.New("failed to reserve room")

func ReserveRoom(ctx context.Context, param ReserveRequest) (db.ReserveResponse, error) {
	reservedRoom, err := db.ReserveRoom(ctx, param)
	if err != nil {
		return db.ReserveRoom{}, ErrReservationFailed
	}
	return reservedRoom, nil

}
