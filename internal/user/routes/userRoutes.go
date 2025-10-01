package routes

import (
	"Travel_Sync/internal/middleware"
	"Travel_Sync/internal/security/config"
	"Travel_Sync/internal/security/service"
	handler "Travel_Sync/internal/user/hander"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(router *gin.Engine, userHandler *handler.UserHandler, jwtService *service.JWTService) {
	api := router.Group("/api")
	{
		user := api.Group("/user")
		// Apply JWT middleware to all user routes
		user.Use(config.JWTMiddleware(jwtService))
		// Apply rate limiting to user routes
		user.Use(middleware.GeneralRateLimiter())
		{
			//user.POST("", userHandler.CreateUser)
			user.DELETE("/:id", userHandler.DeleteUser)
			user.PUT("/:id", userHandler.UpdateUser)
			user.GET("/:id", userHandler.GetUserById)
			user.GET("", userHandler.GetAllUser)
		}
	}
}
