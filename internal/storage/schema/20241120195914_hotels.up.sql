CREATE TABLE public.hotels(
    id uuid PRIMARY KEY DEFAULT gen_random(),
    name VARCHAR(255) NOT NULL,
    star INT NULL,
    location []Float,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
)