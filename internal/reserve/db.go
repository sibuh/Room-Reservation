package reserve

import (
	"context"
	"reservation/internal/storage/db"
)

type Querier interface {
	HoldRoom(ctx context.Context, arg db.HoldRoomParams) (db.Room, error)
}
