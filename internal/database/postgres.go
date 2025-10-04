package database

import (
	"Travel_Sync/internal/config"
	tentity "Travel_Sync/internal/travel/entity"
	"Travel_Sync/internal/user/entity"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(cfg *config.AppConfig) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.PostgresURI), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Set timezone to UTC for all database operations
	if _, err := sqlDB.Exec("SET timezone = 'UTC'"); err != nil {
		log.Printf("Warning: Failed to set timezone to UTC: %v", err)
	}

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	// Automigrate schemas
	if err := db.AutoMigrate(&tentity.TravelTicket{}, &entity.User{}); err != nil {
		return nil, err
	}

	log.Println(" Connected to Postgres via GORM")
	return db, nil
}

func Disconnect(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Println("Failed to get sql.DB for closing:", err)
		return
	}

	if err := sqlDB.Close(); err != nil {
		log.Println("Failed to close database connection:", err)
	} else {
		log.Println("Database connection closed")
	}
}
