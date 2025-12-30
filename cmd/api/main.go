package main

import (
	"log"

	"github.com/abdulwahidkahar/go-booking-reservation-service.git/internal/config"
	"github.com/abdulwahidkahar/go-booking-reservation-service.git/internal/database"
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
}
