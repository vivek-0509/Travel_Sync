package middleware

import (
	"time"

	"Travel_Sync/internal/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupCORS(appCfg *config.AppConfig) gin.HandlerFunc {
	corsCfg := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	if len(appCfg.AllowedOrigins) > 0 {
		corsCfg.AllowOrigins = appCfg.AllowedOrigins
	} else {
		corsCfg.AllowOrigins = []string{"http://localhost:3000", "http://localhost:8080", "http://127.0.0.1:3000", "http://127.0.0.1:8080"}
	}
	return cors.New(corsCfg)
}
