package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewGinRouter() *gin.Engine {
	r := gin.Default()

	// Add health check endpoint (no authentication required)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"message": "Travel Sync API is running",
		})
	})

	return r
}
