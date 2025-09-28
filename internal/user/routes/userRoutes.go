package routes

import (
	handler "Travel_Sync/internal/user/hander"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(router *gin.Engine, userHandler *handler.UserHandler) {
	api := router.Group("/api")
	{
		user := api.Group("/user")
		{
			user.POST("", userHandler.CreateUser)
			user.DELETE("/:id", userHandler.DeleteUser)
			user.PUT("/:id", userHandler.UpdateUser)
			user.GET("/:id", userHandler.GetUserById)
			user.GET("", userHandler.GetAllUser)
		}
	}
}
