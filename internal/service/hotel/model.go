package hotel

import (
	"errors"
	"reservation/pkg/checkzerouuid"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
)

type SearchHotelParam struct {
	Location struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"location"`
}

type RegisterHotelParam struct {
	Name     string    `json:"name"`
	Rating   float64   `json:"rating"`
	OwnerID  uuid.UUID `json:"owner_id"`
	ImageURL string    `json:"image_url"`
	Location struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"location"`
}

func ValidateUUID(value interface{}) error {

	id, ok := value.(uuid.UUID)
	if !ok {
		return errors.New("invalid uuid")
	}

	if checkzerouuid.CheckIfZero(id) {
		return errors.New("owner id is required")
	}
	return nil
}

func (r RegisterHotelParam) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required.Error("name is required")),
		validation.Field(&r.Location, validation.Required.Error("location is required")),
		validation.Field(&r.ImageURL, validation.Required.Error("hotel image is required")),
		validation.Field(&r.OwnerID, validation.Required.Error("owner id is required"),
			validation.By(ValidateUUID)),
		validation.Field(&r.Rating, validation.Required.Error("hotel rating is required")),
	)
}
