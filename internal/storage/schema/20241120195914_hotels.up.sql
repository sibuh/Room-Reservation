
create TYPE hotel_status AS ENUM('PENDING','VERIFIED');

CREATE TABLE public.hotels(
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    owner_id uuid NOT NULL references users(id),
    rating FLOAT NOT NULL DEFAULT 1,
    country VARCHAR(255) NOT NULL,
    city VARCHAR(255) NOT NULL,
    location FLOAT[] NOT NULL,
    image_urls VARCHAR(255)[] NOT NULL,
    status hotel_status NOT NULL DEFAULT 'PENDING',
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);