package main

import (
	"Travel_Sync/internal/config"
	"Travel_Sync/internal/database"
	"Travel_Sync/internal/security/authConfig"
	handler2 "Travel_Sync/internal/security/handler"
	routes2 "Travel_Sync/internal/security/routes"
	securityService "Travel_Sync/internal/security/service"
	"Travel_Sync/internal/server"
	travelHandler "Travel_Sync/internal/travel/handler"
	travelRepo "Travel_Sync/internal/travel/repository"
	travelRoutes "Travel_Sync/internal/travel/routes"
	travelService "Travel_Sync/internal/travel/service"
	handler "Travel_Sync/internal/user/hander"
	"Travel_Sync/internal/user/repository"
	"Travel_Sync/internal/user/routes"
	userService "Travel_Sync/internal/user/service"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading from system env")
	}

	cfg := config.LoadConfig()

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to PostgresDB: %v", err)
	}
	defer database.Disconnect(db)

	// --- Repos & Services ---
	userRepo := repository.NewUserRepo(db)
	userSvc := userService.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userSvc)

	tRepo := travelRepo.NewTravelTicketRepo(db)
	tSvc := travelService.NewTravelTicketService(tRepo, userRepo)
	tHandler := travelHandler.NewTravelTicketHandler(tSvc)

	oauth2Config := authConfig.GetGoogleOAuthConfig()
	authSvc := securityService.NewAuthService(userSvc)
	jwtSvc := securityService.NewJWTService()
	customOAuthSvc := securityService.NewCustomOAuth2Service(oauth2Config, authSvc, jwtSvc)
	authHandler := handler2.NewOAuthHandler(customOAuthSvc)

	// --- Gin Router ---
	ginEngine := server.NewGinRouter()

	// ✅ Optional: Handle preflight requests
	ginEngine.OPTIONS("/*path", func(c *gin.Context) {
		c.Status(200)
	})

	// --- Register routes ---
	routes.RegisterUserRoutes(ginEngine, userHandler, jwtSvc)
	travelRoutes.RegisterTravelRoutes(ginEngine, tHandler, jwtSvc)
	routes2.RegisterAuthRoutes(ginEngine, authHandler, jwtSvc)

	// --- Start server ---
	addr := ":" + cfg.Port
	log.Printf("Listening on %s", addr)
	srv := &http.Server{Addr: addr, Handler: ginEngine}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// --- Graceful shutdown ---
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), server.ShutdownTimeout())
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server exiting")
}
