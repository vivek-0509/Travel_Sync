package database

import (
	"Travel_Sync/internal/config"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"time"
)

type PostgresClient struct {
	Pool *pgxpool.Pool
}

func Connect(cfg config.AppConfig) (*PostgresClient, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, cfg.PostgresURI)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}

	log.Println("Connected to database")

	return &PostgresClient{
		Pool: pool,
	}, nil

}
