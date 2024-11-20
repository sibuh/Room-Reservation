package reserve

import "context"

type db interface {
	ReserveRoom(ctx context.Context, param db.ReserveParam) (db.ReserveRoom, error)
}
