SELECT 
    h.id, 
    h.name, 
    h.owner_id, 
    h.rating, 
    h.country, 
    h.city, 
    h.location, 
    h.image_urls, 
    h.status, 
    h.created_at, 
    h.updated_at, 
    rt.min_price
FROM hotels h
JOIN rooms r ON r.hotel_id = h.id
JOIN (
    SELECT 
        id, 
        MIN(price) AS min_price 
    FROM room_types 
    GROUP BY id
) rt ON rt.id = r.room_type_id
WHERE 
    (h.city LIKE 'addis ababa' OR h.country LIKE 'addis ababa')
    AND h.status = 'VERIFIED'
    AND rt.min_price >= 1
    AND r.id NOT IN (
        SELECT id 
        FROM reservations 
        WHERE 
            (from_time BETWEEN '2025-01-22T06:04:48.943836Z' AND '2025-01-23T06:04:48.943836Z' 
             OR to_time BETWEEN '2025-01-22T06:04:48.943836Z' AND '2025-01-23T06:04:48.943836Z')
            AND reservations.status IN ('SUCCESSFUL', 'PENDING')
    );