package main

import (
	"Travel_Sync/internal/config"
	"Travel_Sync/internal/database"
	"Travel_Sync/internal/server"
	"Travel_Sync/internal/user/hander"
	"Travel_Sync/internal/user/repository"
	"Travel_Sync/internal/user/routes"
	"Travel_Sync/internal/user/service"
	"log"
)

func main() {
	cfg := config.LoadConfig()

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to PostgresDB: %v", err)
	}

	defer database.Disconnect(db) //  ensures DB is closed on exit

	userRepo := repository.NewUserRepo(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	ginEngine := server.NewGinRouter()
	routes.RegisterUserRoutes(ginEngine, userHandler)

	addr := ":" + cfg.Port
	log.Printf("Listening on %s", addr)
	if err := ginEngine.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
