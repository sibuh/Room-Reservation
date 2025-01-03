// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: hotels.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createHotel = `-- name: CreateHotel :one
insert into hotels(name,owner_id,location,rating,image_url)values($1,$2,$3,$4,$5)
 returning id, name, owner_id, rating, location, image_url, status, created_at, updated_at
`

type CreateHotelParams struct {
	Name     string
	OwnerID  pgtype.UUID
	Location []float64
	Rating   float64
	ImageUrl string
}

func (q *Queries) CreateHotel(ctx context.Context, arg CreateHotelParams) (Hotel, error) {
	row := q.db.QueryRow(ctx, createHotel,
		arg.Name,
		arg.OwnerID,
		arg.Location,
		arg.Rating,
		arg.ImageUrl,
	)
	var i Hotel
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.OwnerID,
		&i.Rating,
		&i.Location,
		&i.ImageUrl,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getHotels = `-- name: GetHotels :many
 select id, name, owner_id, rating, location, image_url, status, created_at, updated_at from hotels
`

func (q *Queries) GetHotels(ctx context.Context) ([]Hotel, error) {
	rows, err := q.db.Query(ctx, getHotels)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Hotel
	for rows.Next() {
		var i Hotel
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.OwnerID,
			&i.Rating,
			&i.Location,
			&i.ImageUrl,
			&i.Status,
			&i.CreatedAt,
			&i.UpdatedAt,
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
