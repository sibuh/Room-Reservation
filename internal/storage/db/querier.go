// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type Querier interface {
	CreateHotel(ctx context.Context, arg CreateHotelParams) (Hotel, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	GetUser(ctx context.Context, email string) (User, error)
	GetUserByID(ctx context.Context, id pgtype.UUID) (User, error)
	UpdateRoom(ctx context.Context, arg UpdateRoomParams) (Room, error)
}

var _ Querier = (*Queries)(nil)
