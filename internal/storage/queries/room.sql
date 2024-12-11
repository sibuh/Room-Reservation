-- name: HoldRoom :one
update rooms set status ='HELD',user_id =$1 
where id= $2  
returning *;