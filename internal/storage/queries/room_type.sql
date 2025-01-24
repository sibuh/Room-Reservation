-- name: CreateRoomType :one
INSERT INTO 
room_types (room_type,description,price,capacity) 
VALUES($1,$2,$3,$4)
RETURNING *;