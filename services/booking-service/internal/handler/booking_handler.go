package handler

import (
	"booking-service/constant"
	"booking-service/internal/usecase"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type BookingHandler struct {
	usecase usecase.BookingUsecase
}

func NewBookingHandler(u usecase.BookingUsecase) *BookingHandler {
	return &BookingHandler{u}
}

// POST /bookings
func (h *BookingHandler) CreateBooking(c *gin.Context) {
	var req struct {
		SpaceID   uint   `json:"space_id"`
		StartTime string `json:"start_time"`
		EndTime   string `json:"end_time"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "user not found in context"})
		return
	}
	userID, ok := userIDVal.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid userID type"})
		return
	}

	start, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid.start_time"})
		return
	}
	end, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid.end_time"})
		return
	}

	booking, err := h.usecase.BookSpace(userID, req.SpaceID, start, end)
	if err != nil {
		switch {
		case errors.Is(err, constant.ErrInvalidBookingTime):
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		case errors.Is(err, constant.ErrSpaceNotFound):
			c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "booking.created.successfully", "data": booking})
}

// PUT /bookings/:id/status
func (h *BookingHandler) UpdateBookingStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid.booking_id"})
		return
	}

	var req struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	booking, err := h.usecase.UpdateStatus(uint(id), req.Status)
	if err != nil {
		if strings.Contains(err.Error(), "invalid booking status") {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "booking.updated.successfully", "data": booking})
}

func (h *BookingHandler) GetBookingByID(c *gin.Context) {
	idStr := c.Param("id")
	bookingID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid.booking_id"})
		return
	}

	booking, err := h.usecase.GetBookingByID(uint(bookingID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, booking)
}

// === USER API: GetBookingByUserID ===
func (h *BookingHandler) GetBookingByUserID(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}
	userID, ok := userIDVal.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid userID type"})
		return
	}

	bookings, err := h.usecase.GetBookingByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "bookings.fetched.successfully", "data": bookings})
}

// === ADMIN/MOD API: GetAllBooking ===
func (h *BookingHandler) GetAllBooking(c *gin.Context) {
	// Có thể check role từ middleware: admin/mod mới được gọi
	roleVal, _ := c.Get("role")
	role := roleVal.(string)
	if role != "admin" && role != "moderator" {
		c.JSON(http.StatusForbidden, gin.H{"message": "forbidden"})
		return
	}

	bookings, err := h.usecase.GetAllBooking()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "bookings.fetched.successfully", "data": bookings})
}
