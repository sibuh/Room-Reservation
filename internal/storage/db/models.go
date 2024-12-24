// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"database/sql/driver"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

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

type Hotel struct {
	ID        pgtype.UUID
	Name      string
	Rating    pgtype.Float8
	Location  []float64
	ImageUrl  pgtype.Text
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
	RoomNumber string
	UserID     pgtype.UUID
	HotelID    pgtype.UUID
	Price      float64
	Status     RoomStatus
	CreatedAt  pgtype.Timestamptz
	UpdatedAt  pgtype.Timestamptz
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
