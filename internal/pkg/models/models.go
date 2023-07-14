package models

import (
	"booking/internal/forms"
	"time"
)

type TemplateData struct {
	MapString map[string]string
	MapInt    map[string]int
	MapFloat  map[string]float32
	Data      map[string]interface{}
	CSRFToken string
	Flash     string
	Error     string
	Warning   string
	Form      *forms.Form
}
type User struct {
	ID           int
	FirstName    string
	LastName     string
	Email        string
	PasswordHash string
	AccessLevel  string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Role         string
}
type Room struct {
	ID        int       `json:"id"`
	RoomName  string    `json:"room_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type Restriction struct {
	ID              int
	RistrictionName string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
type Reservation struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	RoomID    int       `json:"room_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Room      Room      `json:"room"`
}
type RoomRestriction struct {
	ID            int
	RoomID        int
	RestrictionID int
	ReservationID int
	StartDate     time.Time
	EndDate       time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Room          Room
	Reservation   Reservation
	Restriction   Restriction
}
type EmailData struct {
	From     string
	To       string
	Subject  string
	Content  string
	Template string
}
type RoomAvailabilityRequest struct {
	StartDate string
	EndDate   string
}
type AddRoomRequest struct {
	RoomNumber string
	RoomType   string
}
type LoginRequest struct {
	Email    string
	Password string
}
