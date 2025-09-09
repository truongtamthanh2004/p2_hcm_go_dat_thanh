package main

import (
	"fmt"
	"log"
	"os"
	_ "payment-service/docs"
	"payment-service/internal/config"
	"payment-service/internal/handler"
	"payment-service/internal/repository"
	"payment-service/internal/router"
	"payment-service/internal/usecase"
)

func main() {
	config.ConnectDB()

	bookingServiceURL := os.Getenv("BOOKING_SERVICE_URL")
	if bookingServiceURL == "" {
		log.Fatal("BOOKING_SERVICE_URL environment variable is not set")
	}

	transactionRepo := repository.NewTransactionRepository(config.DB)
	PaymentUsecase := usecase.NewPaymentUsecase(transactionRepo, config.GetVnpayConfig(), bookingServiceURL)
	paymentHandler := handler.NewPaymentHandler(PaymentUsecase)

	r := router.SetupRouter(paymentHandler)

	port := os.Getenv("PAYMENT_SERVICE_PORT")
	if port == "" {
		port = "8085"
	}

	r.Run(fmt.Sprintf(":%s", port))
}
