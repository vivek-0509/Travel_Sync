package server

import (
	"net/http"
	"time"

	"Travel_Sync/internal/config"
	"Travel_Sync/internal/middleware"

	"github.com/gin-gonic/gin"
)

func NewGinRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	// Set Gin mode via env
	if cfg := config.LoadConfig(); cfg.GinMode != "" {
		gin.SetMode(cfg.GinMode)
	}
	appCfg := config.LoadConfig()
	// Trusted proxies (e.g., AWS Load Balancer)
	if len(appCfg.TrustedProxies) > 0 {
		_ = r.SetTrustedProxies(appCfg.TrustedProxies)
	}

	// Add CORS middleware
	r.Use(middleware.SetupCORS(appCfg))

	// Add global rate limiting
	r.Use(middleware.GeneralRateLimiter())

	// Add health check endpoint (no authentication required)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"message": "Travel Sync API is running",
		})
	})

	// Add CORS test endpoint
	r.GET("/cors-test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"cors_enabled": true,
			"origin":       c.GetHeader("Origin"),
			"message":      "CORS is working correctly",
		})
	})

	return r
}

// ShutdownTimeout returns the time to wait for graceful shutdown
func ShutdownTimeout() time.Duration { return 10 * time.Second }
