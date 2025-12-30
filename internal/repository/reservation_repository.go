package repository

import (
	"database/sql"
	"time"

	"github.com/abdulwahidkahar/go-booking-reservation-service.git/internal/model"
)

type ReservationRepository interface {
	Create(tx *sql.Tx, reservation *model.Reservation) (int64, error)
	GetByID(tx *sql.Tx, reservationID int64) (*model.Reservation, error)
	FindExpiredReservations(tx *sql.Tx, now time.Time) ([]*model.Reservation, error)
	UpdateStatus(tx *sql.Tx, reservationID int64, status model.ReservationStatus) error
}
