-- name: CreateUser :one 
insert into users 
(first_name,last_name,phone_number,email,password,username)
values($1,$2,$3,$4,$5,$6) returning *;
-- name: GetUser :one
select * from users
where email=$1;