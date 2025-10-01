package routes

import (
	"Travel_Sync/internal/middleware"
	"Travel_Sync/internal/security/config"
	secservice "Travel_Sync/internal/security/service"
	thandler "Travel_Sync/internal/travel/handler"

	"github.com/gin-gonic/gin"
)

func RegisterTravelRoutes(router *gin.Engine, handler *thandler.TravelTicketHandler, jwtService *secservice.JWTService) {
	api := router.Group("/api")
	travel := api.Group("/travel")
	travel.Use(config.JWTMiddleware(jwtService))
	{
		travel.POST("", handler.Create)
		travel.GET("", handler.GetAll)
		travel.GET("/my", handler.GetMyTickets)
		travel.GET("/:id", handler.GetByID)
		travel.PUT("/:id", handler.Update)
		travel.DELETE("/:id", handler.Delete)
		travel.GET("/user-responses", handler.GetUserResponses)

		// Apply stricter rate limiting for recommendation endpoint
		recommendations := travel.Group("/")
		recommendations.Use(middleware.RecommendationRateLimiter())
		{
			recommendations.GET("/:id/recommendations", handler.GetRecommendations)
		}
	}
}
