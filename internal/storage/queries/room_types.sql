-- name: AddRoomType :one
insert into room_types(room_type,price,description,max_accupancy)values($1,$2,$3,$4) returning *;