package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

// search room query
const SearchHotelsQuery = `
SELECT 
    h.id, 
    h.name, 
    h.owner_id, 
    h.rating, 
    h.country, 
    h.city, 
    h.location, 
    h.image_urls, 
    h.status, 
    h.created_at, 
    h.updated_at, 
    rt.min_price
FROM hotels h
JOIN rooms r ON r.hotel_id = h.id
JOIN (
    SELECT 
        id, 
        MIN(price) AS min_price 
    FROM room_types 
    GROUP BY id
) rt ON rt.id = r.room_type_id
WHERE 
    (ILIKE '%' || $1 || '%' OR h.country ILIKE '%' || $1 || '%')
    AND h.status = 'VERIFIED'
    AND rt.min_price >= $2
    AND r.id NOT IN (
        SELECT id 
        FROM reservations 
        WHERE 
            (from_time BETWEEN $3 AND $4 
             OR to_time BETWEEN $3 AND $4)
            AND reservations.status IN ('SUCCESSFUL', 'PENDING')
    );`

type SearchHotelsParams struct {
	City     string             `json:"city"`
	Capacity int32              `json:"capacity"`
	FromTime pgtype.Timestamptz `json:"from_time"`
	ToTime   pgtype.Timestamptz `json:"to_time"`
}

type SearchHotelsRow struct {
	ID        pgtype.UUID        `json:"id"`
	Name      string             `json:"name"`
	OwnerID   pgtype.UUID        `json:"owner_id"`
	Rating    float64            `json:"rating"`
	Country   string             `json:"country"`
	City      string             `json:"city"`
	Location  []float64          `json:"location"`
	ImageUrls []string           `json:"image_urls"`
	Status    HotelStatus        `json:"status"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at"`
	MinPrice  float64            `json:"min_price"`
}

func SearchHotels(ctx context.Context, conn *pgxpool.Conn, arg SearchHotelsParams) ([]SearchHotelsRow, error) {
	rows, err := conn.Query(ctx, SearchHotelsQuery,
		arg.City,
		arg.Capacity,
		arg.FromTime,
		arg.ToTime,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SearchHotelsRow
	for rows.Next() {
		var i SearchHotelsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.OwnerID,
			&i.Rating,
			&i.Country,
			&i.City,
			&i.Location,
			&i.ImageUrls,
			&i.Status,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.MinPrice,
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
