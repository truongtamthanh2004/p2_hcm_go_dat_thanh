package router

import (
	"booking-service/internal/handler"
	"booking-service/internal/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(bookingHandler *handler.BookingHandler) *gin.Engine {
	router := gin.Default()
	router.HandleMethodNotAllowed = true // return 405 on wrong method

	router.POST("/api/v1/bookings", middleware.RequireAuth("user"), bookingHandler.CreateBooking)
	router.PUT("/api/v1/bookings/:id/status", bookingHandler.UpdateBookingStatus)
	router.GET("/api/v1/bookings/:id", bookingHandler.GetBookingByID)
	router.GET("/api/v1/bookings/me", middleware.RequireAuth("user"), bookingHandler.GetBookingByUserID)
	router.GET("/api/v1/bookings", middleware.RequireAuth("admin", "moderator"), bookingHandler.GetAllBooking)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
