CREATE TYPE room_status AS ENUM('FREE','HELD','RESERVED','OCCUPAID');

CREATE TABLE public.rooms (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    room_number INT NOT NULL,
    hotel_id uuid NOT NULL references hotels(id),
    room_type_id uuid NOT NULL references room_types(id),
    floor VARCHAR(255) NOT NULL,
    status room_status NOT NULL DEFAULT 'FREE',
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);
