package roomtype

import validation "github.com/go-ozzo/ozzo-validation"

type CreateRoomTypeRequest struct {
	RoomType    string  `json:"room_type"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Capacity    int32   `json:"capacity"`
}

func (c CreateRoomTypeRequest) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.RoomType, validation.Required.Error("room_type is required")),
		validation.Field(&c.Description, validation.Required.Error("description is required")),
		validation.Field(&c.Price, validation.Required.Error("price is required")),
		validation.Field(&c.Capacity, validation.Required.Error("capacity is required")),
	)
}
