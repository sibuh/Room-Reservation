package reserve

import (
	"context"
	"errors"
	"reservation/internal/storage/db"
	"time"

	"github.com/google/uuid"
)

type ReserveRequest struct {
	HotelID uuid.UUID `json:"hotel_id"`
	RoomID  uuid.UUID `json:"room_id"`
	UserID  uuid.UUID `json:"user_id"`
}
type ReserveResponse struct {
	ID         uuid.UUID
	RoomNumber string
	UserID     uuid.UUID
	HotelID    uuid.UUID
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
type Session struct {
	SessionID  uuid.UUID `json:"session_id"`
	PaymentURL string    `json:"payment_url"`
	CancelURL  string    `json:"cancel_url"`
}

var ErrReservationFailed = errors.New("failed to reserve room")

type reserve struct {
	Querier
}

func Init(q Querier) Querier {
	return reserve{
		Querier: q,
	}
}

func (r *reserve) ReserveRoom(ctx context.Context, param ReserveRequest) (ReserveResponse, error) {
	reservedRoom, err := r.Querier.HoldRoom(ctx, db.HoldRoomParams{
		UserID: uuid.NullUUID{
			UUID:  param.UserID,
			Valid: true,
		},
		HotelID: param.HotelID,

		ID: param.RoomID,
	})
	if err != nil {
		return ReserveResponse{}, ErrReservationFailed
	}
	return ReserveResponse{
		ID:         reservedRoom.ID,
		RoomNumber: reservedRoom.RoomNumber,
		UserID:     reservedRoom.UserID.UUID,
		HotelID:    reservedRoom.HotelID,
		CreatedAt:  reservedRoom.CreatedAt,
		UpdatedAt:  reservedRoom.UpdatedAt,
	}, nil

}
