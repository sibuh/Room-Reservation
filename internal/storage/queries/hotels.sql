-- name: CreateHotel :one 
insert into hotels(name,owner_id,location,rating,image_url)values($1,$2,$3,$4,$5)
 returning *;


-- name: GetHotels :many 
 select * from hotels;