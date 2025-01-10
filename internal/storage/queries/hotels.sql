-- name: CreateHotel :one 
insert into hotels(name,owner_id,location,rating,image_urls)values($1,$2,$3,$4,$5)
 returning *;


-- name: GetHotels :many 
select * from hotels limit 10;

-- name: SearchHotels :many
 select h.*,r.*,rt.* from hotels h
 join rooms r on r.hotel_id=h.id
 join room_types rt on rt.room_id=r.id 
 where h.city LIKE $1 or h.country LIKE $1
 and r.id not in(select id from reservations where from_time between $2 and $3 or to_time between $2 and $3)
 ;

-- name: GetHotelByName :one 
 select * from hotels where name=$1;