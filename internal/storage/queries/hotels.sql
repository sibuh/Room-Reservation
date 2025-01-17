-- name: CreateHotel :one 
insert into hotels(name,city,country,owner_id,location,rating,image_urls)values($1,$2,$3,$4,$5,$6,$7)
 returning *;


-- name: GetHotels :many 
select * from hotels limit 10;

-- name: SearchHotels :many
select h.* 
 from hotels h
 join rooms r on r.hotel_id=h.id
 where h.city LIKE $1 or h.country LIKE $1
 and h.status = 'VERIFIED'
 and r.room_type_id in(select id from room_types where capacity >= $2)
 and r.id not in(
    select id from reservations where from_time between $3 and $4
                                                or to_time between $3 and $4
                                                and reservations.status in ( 'SUCCESSFUL','PENDING'));
-- name: GetHotelByName :one 
 select * from hotels where name=$1;