package middleware

import (
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
		AllowAllOrigins:  true, // ✅ allow all origins
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false, // ❌ must be false if AllowAllOrigins=true
		MaxAge:           12 * time.Hour,
	})
}
