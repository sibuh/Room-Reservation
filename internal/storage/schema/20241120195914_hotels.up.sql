CREATE TABLE public.hotels(
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    owner_id uuid references users(id),
    rating Float NULL,
    location Float[],
    image_url VARCHAR(255),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
)