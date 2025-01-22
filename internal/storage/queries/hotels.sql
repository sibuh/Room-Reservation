-- name: CreateHotel :one 
insert into hotels(name,city,country,owner_id,location,rating,image_urls)values($1,$2,$3,$4,$5,$6,$7)
 returning *;


-- name: GetHotels :many 
select * from hotels limit 10;
-- name: GetHotelByName :one 
 select * from hotels where name=$1;