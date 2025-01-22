-- name: CreateHotel :one 
insert into hotels(name,city,country,owner_id,location,rating,image_urls)values($1,$2,$3,$4,$5,$6,$7)
 returning *;


-- name: GetHotels :many 
select * from hotels limit 10;
-- name: GetHotelByName :one 
 select * from hotels where name=$1;

-- name: SearchHotels :many
SELECT DISTINCT h.*, rt.price
FROM hotels h
JOIN rooms r ON h.id = r.hotel_id
JOIN room_types rt ON r.room_type_id = rt.id
LEFT JOIN reservations res 
    ON r.id = res.room_id
    AND res.status IN ('PENDING', 'SUCCESSFUL') -- Only consider active reservations
    AND (
        (res.from_time BETWEEN $1 AND $2) OR (res.to_time  BETWEEN $1 AND  $2 ) -- Overlapping reservation
    )
WHERE h.country = $3 OR h.city = $3
  AND rt.capacity >= $4
  AND res.id IS NULL -- Room is not reserved in the given time range
ORDER BY h.name;
 
-- name: VerifyHotel :one
UPDATE hotels SET status='VERIFIED' WHERE id=$1 RETURNING *;