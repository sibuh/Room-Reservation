// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: room_types.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createRoomType = `-- name: CreateRoomType :one
INSERT INTO 
room_types (room_type,description,price,capacity) 
VALUES($1,$2,$3,$4)
RETURNING id, room_type, price, description, capacity, created_at, updated_at, deleted_at
`

type CreateRoomTypeParams struct {
	RoomType    Roomtype `json:"room_type"`
	Description string   `json:"description"`
	Price       float64  `json:"price"`
	Capacity    int32    `json:"capacity"`
}

func (q *Queries) CreateRoomType(ctx context.Context, arg CreateRoomTypeParams) (RoomType, error) {
	row := q.db.QueryRow(ctx, createRoomType,
		arg.RoomType,
		arg.Description,
		arg.Price,
		arg.Capacity,
	)
	var i RoomType
	err := row.Scan(
		&i.ID,
		&i.RoomType,
		&i.Price,
		&i.Description,
		&i.Capacity,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const getRoomType = `-- name: GetRoomType :one
SELECT id, room_type, price, description, capacity, created_at, updated_at, deleted_at FROM room_types WHERE id=$1
`

func (q *Queries) GetRoomType(ctx context.Context, id pgtype.UUID) (RoomType, error) {
	row := q.db.QueryRow(ctx, getRoomType, id)
	var i RoomType
	err := row.Scan(
		&i.ID,
		&i.RoomType,
		&i.Price,
		&i.Description,
		&i.Capacity,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const getRoomTypes = `-- name: GetRoomTypes :many
SELECT id, room_type, price, description, capacity, created_at, updated_at, deleted_at FROM room_types
`

func (q *Queries) GetRoomTypes(ctx context.Context) ([]RoomType, error) {
	rows, err := q.db.Query(ctx, getRoomTypes)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []RoomType
	for rows.Next() {
		var i RoomType
		if err := rows.Scan(
			&i.ID,
			&i.RoomType,
			&i.Price,
			&i.Description,
			&i.Capacity,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
