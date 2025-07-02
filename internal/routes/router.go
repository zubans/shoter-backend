package routes

import (
	"github.com/gin-gonic/gin"
	"shuter-go/internal/handlers"
	"shuter-go/internal/middlewares"
)

func SetupUserRoutes(r *gin.Engine, userHandler *handlers.UserHandler) {
	api := r.Group("/api/users")
	api.Use(middlewares.LoggingMiddleware())

	{
		api.POST("/create-player-profile", userHandler.Create)
	}
}

//func SetupOrderRoutes(r *gin.Engine, orderHandler *handlers.OrderHandler, authMW gin.HandlerFunc) {
//	orderGroup := r.Group("/api/user/orders")
//	orderGroup.Use(authMW, middlewares.LoggingMiddleware())
//	{
//		orderGroup.POST("", orderHandler.UploadOrder)
//		orderGroup.GET("", orderHandler.GetOrders)
//	}
//}
