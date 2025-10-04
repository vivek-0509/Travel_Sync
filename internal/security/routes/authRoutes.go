package routes

import (
	"Travel_Sync/internal/middleware"
	"Travel_Sync/internal/security/config"
	"Travel_Sync/internal/security/handler"
	"Travel_Sync/internal/security/service"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(router *gin.Engine, authHandler *handler.OAuthHandler, jwtService *service.JWTService) {
	auth := router.Group("/auth")
	// Apply stricter rate limiting for auth endpoints
	auth.Use(middleware.AuthRateLimiter())
	{
		// Public routes (no authentication required)
		auth.GET("/google/login", authHandler.GoogleLogin)
		auth.GET("/google/callback", authHandler.GoogleCallback)

		// Protected routes (authentication required)
		protected := auth.Group("/")
		protected.Use(config.JWTMiddleware(jwtService))
		{
			protected.GET("/me", authHandler.GetCurrentUser)
			protected.POST("/logout", authHandler.Logout)
		}
	}
}
