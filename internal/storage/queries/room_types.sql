-- name: AddRoomType :one
insert into room_types(room_type,price,description,capacity)values($1,$2,$3,$4) returning *;