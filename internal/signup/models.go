package signup

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	PhoneNumber string    `json:"phone_number"`
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type SignupRequest struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	Password    string `json:"password"`
}

func (sr SignupRequest) Validate() error {
	return validation.ValidateStruct(&sr,
		validation.Field(&sr.FirstName, validation.Required.Error("first name is required")),     //add other validation rules
		validation.Field(&sr.LastName, validation.Required.Error("first name is required")),      //add other validation rules
		validation.Field(&sr.PhoneNumber, validation.Required.Error("phone number is required")), //add other validation rules
		validation.Field(&sr.Email, validation.Required.Error("first name is required")),         //add other validation rules
		validation.Field(&sr.Password, validation.Required.Error("password is required")),        //add other validation rules
	)
}

type LoginRequest struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}
