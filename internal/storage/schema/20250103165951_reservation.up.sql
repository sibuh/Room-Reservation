create type reservation_status as enum('PENDING','SUCCESSFUL','FAILED');
create table reservations (
    id uuid primary key default gen_random_uuid(),
    room_id uuid references rooms(id),
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    phone_number VARCHAR(12) NOT NULL,
    email VARCHAR(255) NOT NULL,
    status reservation_status not null default 'PENDING',
    from_time timestamptz NOT NULL,
    to_time timestamptz NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    deleted_at timestamptz NOT NULL
);