-- name: CreateHotel :one 
insert into hotels
 (name,location,rating)values($1,$2,$3)
 returning *;