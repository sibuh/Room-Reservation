SELECT DISTINCT h.id, h.name, h.owner_id, h.rating, h.country, h.city, h.location, h.image_urls, h.status, h.created_at, h.updated_at, rt.capacity
FROM hotels h
JOIN rooms r ON h.id = r.hotel_id
JOIN room_types rt ON r.room_type_id = rt.id
LEFT JOIN reservations res 
    ON r.id = res.room_id
    AND res.status IN ('PENDING', 'SUCCESSFUL') -- Only consider active reservations
    AND (
    
        (res.from_time BETWEEN '2025-01-23T17:52:13.695542Z' AND '2025-01-22T17:52:13.695542Z') OR (res.to_time BETWEEN '2025-01-23T17:52:13.695542Z' AND '2025-01-22T17:52:13.695542Z') -- Overlapping reservation
        )
WHERE h.country = 'addis ababa' OR h.city = 'addis ababa'
  AND rt.capacity >= 2
  AND res.id IS NULL -- Room is not reserved in the given time range
ORDER BY h.name;