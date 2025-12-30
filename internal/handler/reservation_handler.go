package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/abdulwahidkahar/go-booking-reservation-service.git/internal/model"
	"github.com/abdulwahidkahar/go-booking-reservation-service.git/internal/service"
)

type ReservationHandler struct {
	service *service.ReservationService
}

func NewReservationHandler(s *service.ReservationService) *ReservationHandler {
	return &ReservationHandler{service: s}
}

// POST /reserve
func (h *ReservationHandler) ReserveSeat(w http.ResponseWriter, r *http.Request) {
	var req model.ReserveSeatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	res, err := h.service.ReserveSeat(r.Context(), req)
	if err != nil {
		// map DB "no rows" to seat unavailable so we don't leak SQL internals
		if errors.Is(err, sql.ErrNoRows) || err == service.ErrSeatUnavailable {
			h.respondError(w, http.StatusConflict, service.ErrSeatUnavailable.Error())
			return
		}
		// generic internal error (do not return raw SQL errors)
		h.respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	h.respondJSON(w, http.StatusCreated, res)
}

// POST /confirm  (body: { "reservation_id": 123 })
func (h *ReservationHandler) ConfirmPayment(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ReservationID int64 `json:"reservation_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.service.ConfirmPayment(r.Context(), body.ReservationID); err != nil {
		// don't leak sql errors
		if errors.Is(err, sql.ErrNoRows) {
			h.respondError(w, http.StatusNotFound, "reservation not found")
			return
		}
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// POST /expire  (trigger expiration run manually for testing)
func (h *ReservationHandler) ExpireReservations(w http.ResponseWriter, r *http.Request) {
	if err := h.service.ExpireReservations(r.Context()); err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{"status": "expired"})
}

// GET /reservations/{id}
func (h *ReservationHandler) GetReservation(w http.ResponseWriter, r *http.Request) {
	// Expect URL: /reservations/{id}
	idStr := r.URL.Path[len("/reservations/"):] // naive parsing
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid reservation id")
		return
	}

	res, err := h.service.GetReservation(r.Context(), id)
	if err != nil {
		h.respondError(w, http.StatusNotFound, err.Error())
		return
	}
	h.respondJSON(w, http.StatusOK, res)
}

// helpers
func (h *ReservationHandler) respondJSON(w http.ResponseWriter, code int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func (h *ReservationHandler) respondError(w http.ResponseWriter, code int, msg string) {
	h.respondJSON(w, code, map[string]string{"error": msg})
}
