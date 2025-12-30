package repository

import (
	"database/sql"
	"time"

	"github.com/abdulwahidkahar/go-booking-reservation-service.git/internal/model"
)

type SeatRepositoryPG struct {
	db *sql.DB
}

func NewSeatRepositoryPG(db *sql.DB) SeatRepository {
	return &SeatRepositoryPG{db: db}
}

// LockSeat locks a seat row and sets the locked_until timestamp. Returns the seat record.
func (r *SeatRepositoryPG) LockSeat(tx *sql.Tx, seatID int64, lockUntil time.Time) (*model.Seat, error) {
	row := tx.QueryRow("SELECT id, flight_id, seat_number, status, locked_until FROM seats WHERE id = $1 FOR UPDATE", seatID)
	var s model.Seat
	var lockedUntil sql.NullTime
	if err := row.Scan(&s.ID, &s.FlightID, &s.SeatNumber, &s.Status, &lockedUntil); err != nil {
		return nil, err
	}
	if lockedUntil.Valid {
		t := lockedUntil.Time
		s.LockedUntil = &t
	}
	// If seat is available, set it to LOCKED and update locked_until
	if s.Status == model.SeatAvailable {
		_, err := tx.Exec("UPDATE seats SET status = $1, locked_until = $2 WHERE id = $3", model.SeatLocked, lockUntil, seatID)
		if err != nil {
			return nil, err
		}
		t := lockUntil
		s.LockedUntil = &t
		s.Status = model.SeatLocked
	}
	return &s, nil
}

func (r *SeatRepositoryPG) GetSeatByID(tx *sql.Tx, seatID int64) (*model.Seat, error) {
	row := tx.QueryRow("SELECT id, flight_id, seat_number, status, locked_until FROM seats WHERE id = $1", seatID)
	var s model.Seat
	var lockedUntil sql.NullTime
	if err := row.Scan(&s.ID, &s.FlightID, &s.SeatNumber, &s.Status, &lockedUntil); err != nil {
		return nil, err
	}
	if lockedUntil.Valid {
		t := lockedUntil.Time
		s.LockedUntil = &t
	}
	return &s, nil
}

func (r *SeatRepositoryPG) ReleaseExpiredSeats(tx *sql.Tx, now time.Time) error {
	_, err := tx.Exec("UPDATE seats SET status = $1, locked_until = NULL WHERE status = $2 AND locked_until <= $3", model.SeatAvailable, model.SeatLocked, now)
	return err
}

func (r *SeatRepositoryPG) MarkSeatAsBooked(tx *sql.Tx, seatID int64) error {
	_, err := tx.Exec("UPDATE seats SET status = $1, locked_until = NULL WHERE id = $2", model.SeatBooked, seatID)
	return err
}
