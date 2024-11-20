package repository

import (
	"reservation/internal/pkg/models"
	"time"
)

type DatabaseRepo interface {
	MakeReservation(res models.Reservation) (int, error)
	InsertRoomRestriction(restrict models.RoomRestriction) error
	SearchAvailabilityByRoomID(roomID int, startDate, endDate time.Time) (bool, error)
	SearchAvailableRooms(startDate, endDate time.Time) ([]models.Room, error)
	InsertRooms(arg models.AddRoomRequest) error
	Login(arg models.LoginRequest) (string, error)
}
