package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	Port        string
	PostgresURI string
}

func LoadConfig() *AppConfig {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return &AppConfig{
		Port:        os.Getenv("PORT"),
		PostgresURI: os.Getenv("POSTGRES_URI"),
	}

}
