package middleware

import (
	"strings"
	"time"

	"Travel_Sync/internal/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupCORS(appCfg *config.AppConfig) gin.HandlerFunc {
	//corsCfg := cors.Config{
	//	AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	//	AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With", "Access-Control-Request-Method", "Access-Control-Request-Headers"},
	//	ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Credentials"},
	//	AllowCredentials: true,
	//	MaxAge:           12 * time.Hour,
	//}
	//
	//// ✅ Use config origins if provided
	//if len(appCfg.AllowedOrigins) > 0 {
	//	corsCfg.AllowOrigins = appCfg.AllowedOrigins
	//	log.Printf("CORS: Using configured origins: %v", appCfg.AllowedOrigins)
	//} else {
	//	// ✅ Default origins for local + prod
	//	corsCfg.AllowOrigins = []string{
	//		"http://localhost:3000",
	//		"http://localhost:8080",
	//		"http://127.0.0.1:3000",
	//		"http://127.0.0.1:8080",
	//		"https://d3l0cmmj1er9dy.cloudfront.net", // 👈 your frontend (prod)
	//	}
	//	log.Printf("CORS: Using default origins: %v", corsCfg.AllowOrigins)
	//}
	//
	//log.Printf("CORS: Configuration - AllowCredentials: %v, MaxAge: %v", corsCfg.AllowCredentials, corsCfg.MaxAge)
	//
	//return cors.New(corsCfg)

	return cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			// Allow specific origins for production security
			allowedOrigins := []string{
				"http://localhost:3000",
				"http://127.0.0.1:3000",
				"https://www.travelsync.space", // legacy fallback
			}

			// Check if origin is in allowed list
			for _, allowed := range allowedOrigins {
				if origin == allowed {
					return true
				}
			}

			// Allow localhost with any port for development
			if strings.HasPrefix(origin, "http://localhost:") || strings.HasPrefix(origin, "http://127.0.0.1:") {
				return true
			}

			return false
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "Cookie"},
		ExposeHeaders:    []string{"Content-Length", "Set-Cookie"},
		AllowCredentials: true, // ✅ allow cookies
		MaxAge:           12 * time.Hour,
	})
}
