CREATE TABLE public.rooms (
    id uuid PRIMARY KEY DEFAULT gen_random(),
    room_name VARCHAR(255) NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL now()
);
