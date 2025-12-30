package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/abdulwahidkahar/go-booking-reservation-service.git/internal/model"
	"github.com/abdulwahidkahar/go-booking-reservation-service.git/internal/repository"
)

var (
	ErrSeatUnavailable     = errors.New("seat is not available")
	ErrReservationNotFound = errors.New("reservation not found")
)

type ReservationService struct {
	db              *sql.DB
	seatRepo        repository.SeatRepository
	reservationRepo repository.ReservationRepository
}

func NewReservationService(
	db *sql.DB,
	seatRepo repository.SeatRepository,
	reservationRepo repository.ReservationRepository,
) *ReservationService {
	return &ReservationService{
		db:              db,
		seatRepo:        seatRepo,
		reservationRepo: reservationRepo,
	}
}

func (s *ReservationService) ReserveSeat(
	ctx context.Context,
	req model.ReserveSeatRequest,
) (*model.Reservation, error) {

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	lockUntil := time.Now().Add(10 * time.Minute)

	seat, err := s.seatRepo.LockSeat(
		tx,
		req.SeatID,
		lockUntil,
	)
	if err != nil {
		// If the seat row doesn't exist, translate to seat unavailable
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrSeatUnavailable
		}
		return nil, err
	}

	if seat.Status == model.SeatUnavailable {
		return nil, ErrSeatUnavailable
	}

	reservation := &model.Reservation{
		UserID:     req.UserID,
		FlightID:   req.FlightID,
		SeatID:     req.SeatID,
		Status:     model.ReservationPending,
		ExpiryTime: &lockUntil,
	}

	id, err := s.reservationRepo.Create(tx, reservation)
	if err != nil {
		return nil, err
	}
	reservation.ID = id

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return reservation, nil
}

func (s *ReservationService) ConfirmPayment(
	ctx context.Context,
	reservationID int64,
) error {

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	reservation, err := s.reservationRepo.GetByID(tx, reservationID)
	if err != nil {
		return err
	}

	if reservation.Status != model.ReservationPending {
		return errors.New("invalid reservation status")
	}

	if err = s.seatRepo.MarkSeatAsBooked(tx, reservation.SeatID); err != nil {
		return err
	}

	if err = s.reservationRepo.UpdateStatus(tx, reservationID, model.ReservationConfirmed); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *ReservationService) GetReservation(ctx context.Context, reservationID int64) (*model.Reservation, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	res, err := s.reservationRepo.GetByID(tx, reservationID)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return res, nil
}

func (s *ReservationService) ExpireReservations(
	ctx context.Context,
) error {

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	now := time.Now()

	if err = s.seatRepo.ReleaseExpiredSeats(tx, now); err != nil {
		return err
	}

	expiredReservations, err := s.reservationRepo.FindExpiredReservations(tx, now)
	if err != nil {
		return err
	}
	for _, r := range expiredReservations {
		if err = s.reservationRepo.UpdateStatus(tx, r.ID, model.ReservationExpired); err != nil {
			return err
		}
	}

	return tx.Commit()
}
