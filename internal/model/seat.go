package model

import "time"

type SeatStatus string

const (
	SeatUnavailable SeatStatus = "UNAVAILABLE"
	SeatAvailable   SeatStatus = "AVAILABLE"
	SeatLocked      SeatStatus = "LOCKED"
	SeatBooked      SeatStatus = "BOOKED"
)

type Seat struct {
	ID          int64
	FlightID    int64
	SeatNumber  string
	Status      SeatStatus
	LockedUntil *time.Time
}
