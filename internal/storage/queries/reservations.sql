-- name: CreateReservation :one
insert into reservations (room_id,user_id,status,from_time,to_time)
values($1,$2,$3,$4,$5) 
returning *;