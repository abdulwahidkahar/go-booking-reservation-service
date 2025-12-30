package model

import "time"

type ReservationStatus string

const (
	ReservationPending   ReservationStatus = "PENDING"
	ReservationConfirmed ReservationStatus = "CONFIRMED"
	ReservationExpired   ReservationStatus = "EXPIRED"
	ReservationCancelled ReservationStatus = "CANCELLED"
)

type ReserveSeatRequest struct {
	UserID   int64 `json:"user_id"`
	FlightID int64 `json:"flight_id"`
	SeatID   int64 `json:"seat_id"`
}

type Reservation struct {
	ID         int64             `json:"id"`
	UserID     int64             `json:"user_id"`
	FlightID   int64             `json:"flight_id"`
	SeatID     int64             `json:"seat_id"`
	Status     ReservationStatus `json:"status"`
	ExpiryTime *time.Time        `json:"expiry_time"`
}
