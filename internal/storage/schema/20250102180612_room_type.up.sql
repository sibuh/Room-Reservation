CREATE TYPE roomtype AS 
ENUM('SINGLE_ROOM','DOUBLE_ROOM','TWIN_ROOM','TRIPLE_ROOM','QUAD_ROOM','QEEN_ROOM','KING_ROOM','SUIT_ROOM');

CREATE TABLE room_types(
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    room_type roomtype NOT NULL DEFAULT 'SINGLE_ROOM',
    price FLOAT NOT NULL DEFAULT 0.0,
    description STRING NOT NULL,
    max_accupancy  INT NOT NULL ,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    deleted_at timestamptz NULL
);