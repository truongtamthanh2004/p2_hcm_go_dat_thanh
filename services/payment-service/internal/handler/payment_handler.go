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

// CreatePaymentUrl godoc
// @Summary      Create Payment URL
// @Description  Generate a VNPAY payment URL for a booking
// @Tags         payments
// @Accept       json
// @Produce      json
// @Param        booking_id query int true "Booking ID"
// @Success      200 {object} map[string]string "payment_url"
// @Failure      400 {object} map[string]string "error"
// @Router       /payments/create-url [get]
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

// VnpayReturn godoc
// @Summary      VNPAY Return
// @Description  Handle VNPAY return callback after payment
// @Tags         payments
// @Accept       json
// @Produce      json
// @Param        vnp_Amount query string true "Amount"
// @Param        vnp_TxnRef query string true "Transaction Ref"
// @Param        vnp_ResponseCode query string true "Response Code"
// @Param        vnp_SecureHash query string true "Secure Hash"
// @Success      302 {string} string "Redirect to success/failure page"
// @Failure      400 {object} map[string]string "error"
// @Router       /payments/vnpay-return [get]
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
