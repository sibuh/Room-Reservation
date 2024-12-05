-- name: CreateUser :one 
insert into users 
(first_name,last_name,phone_number,email,password)
values($1,$2,$3,$4,$5) returning *;