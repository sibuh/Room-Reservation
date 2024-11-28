package reserve

import (
	"context"
	"errors"
	"reservation/internal/storage/db"

	"github.com/google/uuid"
	"github.com/jackc/pgtype"
)

type ReserveRequest struct {
	HotelID uuid.UUID `json:"hotel_id"`
	RoomID  uuid.UUID `json:"room_id"`
	UserID  uuid.UUID `json:"user_id"`
}
type ReserveResponse struct {
	ID         pgtype.UUID
	RoomNumber string
	UserID     pgtype.UUID
	HotelID    pgtype.UUID
	CreatedAt  pgtype.Timestamptz
	UpdatedAt  pgtype.Timestamptz
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
		UserID:  pgtype.UUID(param.UserID),
		HotelID: pgtype.UUID(param.HotelID),
		ID:      pgtype.UUID(param.RoomID),
	})
	if err != nil {
		return ReserveResponse{}, ErrReservationFailed
	}
	return ReserveResponse{
		ID:         reservedRoom.ID,
		RoomNumber: reservedRoom.RoomNumber,
		UserID:     reservedRoom.UserID,
		HotelID:    reservedRoom.HotelID,
		CreatedAt:  reservedRoom.CreatedAt,
		UpdatedAt:  reservedRoom.UpdatedAt,
	}, nil

}
