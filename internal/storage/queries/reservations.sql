-- name: CreateReservation :one
insert into reservations (room_id,first_name,last_name,phone_number,email,status,from_time,to_time)
values($1,$2,$3,$4,$5,$6,$7,$8) 
returning *;

-- name: UpdateReservation :one 
update reservations set status=$1 where id =$2 returning *;
-- name: GetReservationStatus :one
select status from reservations where id =$1;

-- name: GetRoomReservations :many
select * from reservations where room_id =$1 and (from_time > now() or to_time > now());
-- name: CheckOverlap :one
select count(id) from reservations where room_id=$1 
                                    and (from_time between $2 and $3 or to_time between $2 and $3);