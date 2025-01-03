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
	HotelStatus HotelStatus
	Valid       bool // Valid is true if HotelStatus is not NULL
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
	ReservationStatus ReservationStatus
	Valid             bool // Valid is true if ReservationStatus is not NULL
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
	RoomStatus RoomStatus
	Valid      bool // Valid is true if RoomStatus is not NULL
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
	Roomtype Roomtype
	Valid    bool // Valid is true if Roomtype is not NULL
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
	ID        pgtype.UUID
	Name      string
	OwnerID   pgtype.UUID
	Rating    float64
	Location  []float64
	ImageUrls []string
	Status    HotelStatus
	CreatedAt pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
}

type Reservation struct {
	ID       pgtype.UUID
	RoomID   pgtype.UUID
	UserID   pgtype.UUID
	Status   ReservationStatus
	FromTime pgtype.Timestamptz
	ToTime   pgtype.Timestamptz
}

type Room struct {
	ID         pgtype.UUID
	RoomNumber int32
	HotelID    pgtype.UUID
	RoomTypeID pgtype.UUID
	Floor      string
	Status     RoomStatus
	CreatedAt  pgtype.Timestamptz
	UpdatedAt  pgtype.Timestamptz
}

type RoomType struct {
	ID          pgtype.UUID
	RoomType    Roomtype
	Price       float64
	Description string
	CreatedAt   pgtype.Timestamptz
	UpdatedAt   pgtype.Timestamptz
	DeletedAt   pgtype.Timestamptz
}

type User struct {
	ID          pgtype.UUID
	FirstName   string
	LastName    string
	PhoneNumber string
	Email       string
	Password    string
	Username    string
	CreatedAt   pgtype.Timestamptz
	UpdatedAt   pgtype.Timestamptz
}
