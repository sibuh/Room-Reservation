-- name: CreateReservation :one
insert into reservations (room_id,user_id,status,from_time,to_time)
values($1,$2,$3,$4,$5) 
returning *;

-- name: UpdateReservation :one 
update reservations set status=$1 where id =$2 returning *;

-- name: GetRoomReservations :many
select * from reservations where room_id =$1 and (from_time > now() or to_time > now());