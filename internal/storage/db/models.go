// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"database/sql/driver"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

type HotelStatus string

const (
	HotelStatusPENDING  HotelStatus = "PENDING"
	HotelStatusVERIFIED HotelStatus = "VERIFIED"
)

func (e *HotelStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = HotelStatus(s)
	case string:
		*e = HotelStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for HotelStatus: %T", src)
	}
	return nil
}

type NullHotelStatus struct {
	HotelStatus HotelStatus `json:"hotel_status"`
	Valid       bool        `json:"valid"` // Valid is true if HotelStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullHotelStatus) Scan(value interface{}) error {
	if value == nil {
		ns.HotelStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.HotelStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullHotelStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.HotelStatus), nil
}

type ReservationStatus string

const (
	ReservationStatusPENDING    ReservationStatus = "PENDING"
	ReservationStatusSUCCESSFUL ReservationStatus = "SUCCESSFUL"
	ReservationStatusCANCELLED  ReservationStatus = "CANCELLED"
	ReservationStatusFAILED     ReservationStatus = "FAILED"
)

func (e *ReservationStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = ReservationStatus(s)
	case string:
		*e = ReservationStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for ReservationStatus: %T", src)
	}
	return nil
}

type NullReservationStatus struct {
	ReservationStatus ReservationStatus `json:"reservation_status"`
	Valid             bool              `json:"valid"` // Valid is true if ReservationStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullReservationStatus) Scan(value interface{}) error {
	if value == nil {
		ns.ReservationStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.ReservationStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullReservationStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.ReservationStatus), nil
}

type RoomStatus string

const (
	RoomStatusFREE     RoomStatus = "FREE"
	RoomStatusHELD     RoomStatus = "HELD"
	RoomStatusRESERVED RoomStatus = "RESERVED"
	RoomStatusOCCUPAID RoomStatus = "OCCUPAID"
)

func (e *RoomStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = RoomStatus(s)
	case string:
		*e = RoomStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for RoomStatus: %T", src)
	}
	return nil
}

type NullRoomStatus struct {
	RoomStatus RoomStatus `json:"room_status"`
	Valid      bool       `json:"valid"` // Valid is true if RoomStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullRoomStatus) Scan(value interface{}) error {
	if value == nil {
		ns.RoomStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.RoomStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullRoomStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.RoomStatus), nil
}

type Roomtype string

const (
	RoomtypeSINGLEROOM Roomtype = "SINGLE_ROOM"
	RoomtypeDOUBLEROOM Roomtype = "DOUBLE_ROOM"
	RoomtypeTWINROOM   Roomtype = "TWIN_ROOM"
	RoomtypeTRIPLEROOM Roomtype = "TRIPLE_ROOM"
	RoomtypeQUADROOM   Roomtype = "QUAD_ROOM"
	RoomtypeQEENROOM   Roomtype = "QEEN_ROOM"
	RoomtypeKINGROOM   Roomtype = "KING_ROOM"
	RoomtypeSUITROOM   Roomtype = "SUIT_ROOM"
)

func (e *Roomtype) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = Roomtype(s)
	case string:
		*e = Roomtype(s)
	default:
		return fmt.Errorf("unsupported scan type for Roomtype: %T", src)
	}
	return nil
}

type NullRoomtype struct {
	Roomtype Roomtype `json:"roomtype"`
	Valid    bool     `json:"valid"` // Valid is true if Roomtype is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullRoomtype) Scan(value interface{}) error {
	if value == nil {
		ns.Roomtype, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.Roomtype.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullRoomtype) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.Roomtype), nil
}

type Hotel struct {
	ID        pgtype.UUID        `json:"id"`
	Name      string             `json:"name"`
	OwnerID   pgtype.UUID        `json:"owner_id"`
	Rating    float64            `json:"rating"`
	Country   string             `json:"country"`
	City      string             `json:"city"`
	Location  []float64          `json:"location"`
	ImageUrls []string           `json:"image_urls"`
	Status    HotelStatus        `json:"status"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at"`
}

type Reservation struct {
	ID          pgtype.UUID        `json:"id"`
	RoomID      pgtype.UUID        `json:"room_id"`
	FirstName   string             `json:"first_name"`
	LastName    string             `json:"last_name"`
	PhoneNumber string             `json:"phone_number"`
	Email       string             `json:"email"`
	Status      ReservationStatus  `json:"status"`
	FromTime    pgtype.Timestamptz `json:"from_time"`
	ToTime      pgtype.Timestamptz `json:"to_time"`
	CreatedAt   pgtype.Timestamptz `json:"created_at"`
	UpdatedAt   pgtype.Timestamptz `json:"updated_at"`
	DeletedAt   pgtype.Timestamptz `json:"deleted_at"`
}

type Room struct {
	ID         pgtype.UUID        `json:"id"`
	RoomNumber int32              `json:"room_number"`
	HotelID    pgtype.UUID        `json:"hotel_id"`
	RoomTypeID pgtype.UUID        `json:"room_type_id"`
	Floor      string             `json:"floor"`
	Status     RoomStatus         `json:"status"`
	CreatedAt  pgtype.Timestamptz `json:"created_at"`
	UpdatedAt  pgtype.Timestamptz `json:"updated_at"`
}

type RoomType struct {
	ID           pgtype.UUID        `json:"id"`
	RoomType     Roomtype           `json:"room_type"`
	Price        float64            `json:"price"`
	Description  string             `json:"description"`
	MaxAccupancy int32              `json:"max_accupancy"`
	CreatedAt    pgtype.Timestamptz `json:"created_at"`
	UpdatedAt    pgtype.Timestamptz `json:"updated_at"`
	DeletedAt    pgtype.Timestamptz `json:"deleted_at"`
}

type User struct {
	ID          pgtype.UUID        `json:"id"`
	FirstName   string             `json:"first_name"`
	LastName    string             `json:"last_name"`
	PhoneNumber string             `json:"phone_number"`
	Email       string             `json:"email"`
	Password    string             `json:"password"`
	Username    string             `json:"username"`
	CreatedAt   pgtype.Timestamptz `json:"created_at"`
	UpdatedAt   pgtype.Timestamptz `json:"updated_at"`
}
