package handler

import (
	"net/http"
	"payment-service/internal/config"
	"payment-service/internal/usecase"
	"payment-service/utils"

	"strconv"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	usecase usecase.PaymentUsecase
}

func NewPaymentHandler(u usecase.PaymentUsecase) *PaymentHandler {
	return &PaymentHandler{u}
}

func (h *PaymentHandler) CreatePaymentUrl(c *gin.Context) {
	bookingIDStr := c.Query("booking_id")
	bookingID, err := strconv.Atoi(bookingIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid.booking_id"})
		return
	}
	clientIP := c.ClientIP()

	url, err := h.usecase.CreatePaymentUrl(uint(bookingID), clientIP)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"payment_url": url})
}

func (h *PaymentHandler) VnpayReturn(c *gin.Context) {
	if !utils.VerifyVnpSignature(c.Request.URL.Query(), config.GetVnpayConfig().HashSecret) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid.signature"})
		return
	}

	redirectURL, err := h.usecase.HandleVnpReturn(c.Request.URL.Query())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, redirectURL)
}
