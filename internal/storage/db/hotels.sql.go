// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: hotels.sql

package db

import (
	"context"
	"database/sql"
)

const createHotel = `-- name: CreateHotel :one
insert into hotels
 (name,location,rating)values($1,$2,$3)
 returning id, name, rating, location, image_url, created_at, updated_at
`

type CreateHotelParams struct {
	Name     string
	Location []float64
	Rating   sql.NullFloat64
}

func (q *Queries) CreateHotel(ctx context.Context, arg CreateHotelParams) (Hotel, error) {
	row := q.db.QueryRow(ctx, createHotel, arg.Name, arg.Location, arg.Rating)
	var i Hotel
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Rating,
		&i.Location,
		&i.ImageUrl,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
