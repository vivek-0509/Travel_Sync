package config

import (
	"Travel_Sync/internal/security/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// JWTMiddleware creates a middleware for JWT authentication
func JWTMiddleware(jwtService *service.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {

		if c.Request.Method == http.MethodOptions {
			c.Next()
			return
		}

		// Get JWT token from cookie
		tokenString, err := c.Cookie("jwt_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "JWT token not found in cookies"})
			c.Abort()
			return
		}

		// Validate JWT token
		claims, err := jwtService.ValidateJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid JWT token"})
			c.Abort()
			return
		}

		// Set user information in context for use in handlers
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("jwt_claims", claims)

		c.Next()
	}
}
