-- name: UpdateRoom :one
update rooms set status =$1
where id= $2  
returning *;

-- name: GetRoom :one 
select
r.id, r.room_number, 
r.hotel_id, r.room_type_id, 
r.floor, r.status, 
r.created_at, r.updated_at,
rt.price
from rooms r 
join room_types rt on r.room_type_id=rt.id
where r.id =$1;

-- name: SearchRoom :many
select * from rooms 
where price < $1 
and hotel_id in (select id from hotels where ST_DWithin(location, ST_GeogPoint($2, $3), 1000))
and roon_id not in(select id from reservations where from_time between $4 and $5 or to_time between $4 and $5)
and $6 =(select room_type from room_types where id=room_type_id);