package repository

import (
	"booking/internal/pkg/models"
	"time"
)

type DatabaseRepo interface {
	MakeReservation(res models.Reservation) (int, error)
	InsertRoomRestriction(restrict models.RoomRestriction) error
	SearchAvailabilityByRoomID(roomID int, startDate, endDate time.Time) (bool, error)
	SearchAvailableRooms(startDate, endDate time.Time) ([]models.Room, error)
}
