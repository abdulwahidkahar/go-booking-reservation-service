package repository

import (
	"database/sql"
	"time"

	"github.com/abdulwahidkahar/go-booking-reservation-service.git/internal/model"
)

type SeatRepository interface {
	LockSeat(tx *sql.Tx, seatID int64, lockUntil time.Time) (*model.Seat, error)
	GetSeatByID(tx *sql.Tx, seatID int64) (*model.Seat, error)
	ReleaseExpiredSeats(tx *sql.Tx, now time.Time) error
	MarkSeatAsBooked(tx *sql.Tx, seatID int64) error
}
