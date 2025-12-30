package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/abdulwahidkahar/go-booking-reservation-service.git/internal/config"
	"github.com/abdulwahidkahar/go-booking-reservation-service.git/internal/database"
	"github.com/abdulwahidkahar/go-booking-reservation-service.git/internal/handler"
	"github.com/abdulwahidkahar/go-booking-reservation-service.git/internal/repository"
	"github.com/abdulwahidkahar/go-booking-reservation-service.git/internal/service"
	"github.com/joho/godotenv"
)

func main() {

	_ = godotenv.Load()

	cfg := config.LoadConfig()

	db, err := database.NewPostgresDB(
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBSSL,
	)

	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	defer db.Close()

	log.Println("database connected successfully")

	// Repositories
	seatRepo := repository.NewSeatRepositoryPG(db)
	reservationRepo := repository.NewReservationRepositoryPG(db)

	// Services
	resService := service.NewReservationService(db, seatRepo, reservationRepo)

	// Handlers and routes
	h := handler.NewReservationHandler(resService)
	mux := http.NewServeMux()
	mux.HandleFunc("/reserve", h.ReserveSeat)
	mux.HandleFunc("/confirm", h.ConfirmPayment)
	mux.HandleFunc("/expire", h.ExpireReservations)
	mux.HandleFunc("/reservations/", h.GetReservation) // GET /reservations/{id}

	port := cfg.Port
	if port == "" {
		port = "8080"
	}
	addr := fmt.Sprintf(":" + port)

	log.Printf("starting server on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
