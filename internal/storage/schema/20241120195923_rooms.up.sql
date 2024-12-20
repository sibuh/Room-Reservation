CREATE TYPE room_status AS ENUM('FREE','HELD','RESERVED');

CREATE TABLE public.rooms (
    id uuid PRIMARY KEY DEFAULT gen_random(),
    room_number VARCHAR(255) NOT NULL DEFAULT 'G00',
    user_id uuid NULL references users(id),
    hotel_id uuid NOT NULL references hotels(id),
    status room_status NOT NULL DEFAULT 'FREE',
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);
