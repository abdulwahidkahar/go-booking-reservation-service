package repository

import (
	"database/sql"
	"time"

	"github.com/abdulwahidkahar/go-booking-reservation-service.git/internal/model"
)

type ReservationRepositoryPG struct {
	db *sql.DB
}

func NewReservationRepositoryPG(db *sql.DB) ReservationRepository {
	return &ReservationRepositoryPG{db: db}
}

func (r *ReservationRepositoryPG) Create(tx *sql.Tx, reservation *model.Reservation) (int64, error) {
	var id int64
	row := tx.QueryRow(
		"INSERT INTO reservations (user_id, flight_id, seat_id, status, expires_at) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		reservation.UserID,
		reservation.FlightID,
		reservation.SeatID,
		reservation.Status,
		reservation.ExpiryTime,
	)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *ReservationRepositoryPG) GetByID(tx *sql.Tx, reservationID int64) (*model.Reservation, error) {
	row := tx.QueryRow("SELECT id, user_id, flight_id, seat_id, status, expires_at FROM reservations WHERE id = $1", reservationID)
	var res model.Reservation
	var expires sql.NullTime
	if err := row.Scan(&res.ID, &res.UserID, &res.FlightID, &res.SeatID, &res.Status, &expires); err != nil {
		return nil, err
	}
	if expires.Valid {
		t := expires.Time
		res.ExpiryTime = &t
	}
	return &res, nil
}

func (r *ReservationRepositoryPG) FindExpiredReservations(tx *sql.Tx, now time.Time) ([]*model.Reservation, error) {
	rows, err := tx.Query("SELECT id, user_id, flight_id, seat_id, status, expires_at FROM reservations WHERE status = $1 AND expires_at <= $2", model.ReservationPending, now)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*model.Reservation
	for rows.Next() {
		var res model.Reservation
		var expires sql.NullTime
		if err := rows.Scan(&res.ID, &res.UserID, &res.FlightID, &res.SeatID, &res.Status, &expires); err != nil {
			return nil, err
		}
		if expires.Valid {
			t := expires.Time
			res.ExpiryTime = &t
		}
		out = append(out, &res)
	}
	return out, rows.Err()
}

func (r *ReservationRepositoryPG) UpdateStatus(tx *sql.Tx, reservationID int64, status model.ReservationStatus) error {
	_, err := tx.Exec("UPDATE reservations SET status = $1 WHERE id = $2", status, reservationID)
	return err
}
