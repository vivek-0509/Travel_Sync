package config

import (
	"os"
)

type AppConfig struct {
	Port        string
	PostgresURI string
}

func LoadConfig() *AppConfig {

	return &AppConfig{
		Port:        os.Getenv("PORT"),
		PostgresURI: os.Getenv("POSTGRES_URI"),
	}

}
