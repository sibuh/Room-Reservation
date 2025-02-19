-- name: CreateRoomType :one
INSERT INTO 
room_types (room_type,description,price,capacity) 
VALUES($1,$2,$3,$4)
RETURNING *;

-- name: GetRoomType :one
SELECT * FROM room_types WHERE id=$1;

-- name: GetRoomTypes :many
SELECT * FROM room_types;