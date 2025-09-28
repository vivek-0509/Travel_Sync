package database

import (
	"Travel_Sync/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
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

	// Test connection
	if err := sqlDB.Ping(); err != nil {
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
