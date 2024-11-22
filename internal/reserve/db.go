package reserve

import "context"

type Querier interface {
	ReserveRoom(ctx context.Context, param db.ReserveParam) (db.ReserveRoom, error)
}
