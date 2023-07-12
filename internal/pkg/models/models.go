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
	ID          int
	FirstName   string
	LastName    string
	Email       string
	Password    string
	AccessLevel string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
type Room struct {
	ID        int
	RoomName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
type Restriction struct {
	ID              int
	RistrictionName string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
type Reservation struct {
	ID        int
	FirstName string
	LastName  string
	Email     string
	Phone     string
	StartDate time.Time
	EndDate   time.Time
	RoomID    int
	CreatedAt time.Time
	UpdatedAt time.Time
	Room      Room
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
