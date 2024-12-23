create type reservation_status as enum('PENDING','SUCCESSFUL','FAILED');
create table reservations (
    id uuid primary key default gen_random_uuid(),
    room_id uuid references rooms(id),
    user_id uuid references users(id),
    status reservation_status not null default 'PENDING',
    from_time timestamptz not null,
    to_time timestamptz not null
);