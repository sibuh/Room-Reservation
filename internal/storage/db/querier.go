// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type Querier interface {
	AddRoom(ctx context.Context, arg AddRoomParams) (Room, error)
	CheckOverlap(ctx context.Context, arg CheckOverlapParams) (int64, error)
	CreateHotel(ctx context.Context, arg CreateHotelParams) (Hotel, error)
	CreateReservation(ctx context.Context, arg CreateReservationParams) (Reservation, error)
	CreateRoomType(ctx context.Context, arg CreateRoomTypeParams) (RoomType, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	GetHotelByName(ctx context.Context, name string) (Hotel, error)
	GetHotelRooms(ctx context.Context, hotelID pgtype.UUID) ([]GetHotelRoomsRow, error)
	GetHotels(ctx context.Context) ([]Hotel, error)
	GetReservationStatus(ctx context.Context, id pgtype.UUID) (ReservationStatus, error)
	GetRoom(ctx context.Context, id pgtype.UUID) (GetRoomRow, error)
	GetRoomReservations(ctx context.Context, roomID pgtype.UUID) ([]Reservation, error)
	GetRoomType(ctx context.Context, id pgtype.UUID) (RoomType, error)
	GetUser(ctx context.Context, email string) (User, error)
	GetUserByID(ctx context.Context, id pgtype.UUID) (User, error)
	SearchHotels(ctx context.Context, arg SearchHotelsParams) ([]SearchHotelsRow, error)
	SearchRoom(ctx context.Context, arg SearchRoomParams) ([]Room, error)
	UpdateReservation(ctx context.Context, arg UpdateReservationParams) (Reservation, error)
	UpdateRoom(ctx context.Context, arg UpdateRoomParams) (Room, error)
	VerifyHotel(ctx context.Context, id pgtype.UUID) (Hotel, error)
}

var _ Querier = (*Queries)(nil)
