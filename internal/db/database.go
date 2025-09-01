package db

import (
	"fleet-pulse-users-service/internal/config"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func DatabaseConnection() *gorm.DB {
	settings := config.Get()
	sqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		settings.Database.Host,
		settings.Database.Port,
		settings.Database.User,
		settings.Database.Password,
		settings.Database.Name,
	)
	log.Printf("Connecting to database")
	db, err := gorm.Open(postgres.Open(sqlInfo), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database %s", err)
	}
	log.Printf("Successfully connected to database")
	return db
}
