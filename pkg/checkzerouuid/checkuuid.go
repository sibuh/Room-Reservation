package checkzerouuid

import "github.com/google/uuid"

const zerozUUID = "00000000-0000-0000-0000-000000000000"

func CheckIfZero(id uuid.UUID) bool {
	return id.String() == zerozUUID
}
