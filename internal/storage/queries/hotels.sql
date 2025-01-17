-- name: CreateHotel :one 
insert into hotels(name,city,country,owner_id,location,rating,image_urls)values($1,$2,$3,$4,$5,$6,$7)
 returning *;


-- name: GetHotels :many 
select * from hotels limit 10;

-- name: SearchHotels :many

with cte as (
select 
h.id hid, 
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
r.id rid, 
r.room_number, 
r.hotel_id, 
r.room_type_id, 
r.floor, 
r.status, 
r.created_at, 
r.updated_at from hotels h
 join rooms r on r.hotel_id=h.id 
 where h.city LIKE $1 or h.country LIKE $1
 and r.id not in(
    select id from reservations where from_time between $2 and $3 
                                                or to_time between $2 and $3
                                                and reservations.status in ( 'SUCCESSFUL','PENDING'))
)
select rt.*,c.* from room_types rt
join cte c on c.room_type_id=rt.id;



-- name: GetHotelByName :one 
 select * from hotels where name=$1;