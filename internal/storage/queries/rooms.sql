-- name: UpdateRoom :one
update rooms set status =$1,user_id =$2 
where id= $3  
returning *;

-- name: GetRoom :one 
select * from rooms where id =$1;