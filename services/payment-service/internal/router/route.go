package router

import (
	"payment-service/internal/handler"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(paymentHandler *handler.PaymentHandler) *gin.Engine {
	router := gin.Default()
	router.HandleMethodNotAllowed = true // return 405 on wrong method

	paymentGroup := router.Group("/api/v1/payments")
	{
		// GET: /api/payments/create?booking_id=123
		paymentGroup.GET("/create", paymentHandler.CreatePaymentUrl)

		// GET: /api/payments/vnpay/callback?...
		paymentGroup.GET("/vnpay/callback", paymentHandler.VnpayReturn)
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
