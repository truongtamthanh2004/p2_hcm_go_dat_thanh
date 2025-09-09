package handler

import (
	"booking-service/constant"
	"booking-service/internal/dto"
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

type BookingRequest struct {
	SpaceID   uint   `json:"space_id"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

type UpdateBookingStatusRequest struct {
	Status string `json:"status"`
}

// CreateBooking godoc
// @Summary      Create a booking
// @Description  User tạo booking mới cho 1 space
// @Tags         bookings
// @Accept       json
// @Produce      json
// @Param        request body BookingRequest true "Booking request"
// @Success      201 {object} map[string]interface{} "booking created"
// @Failure      400 {object} map[string]string "invalid input"
// @Failure      401 {object} map[string]string "unauthorized"
// @Failure      404 {object} map[string]string "space not found"
// @Router       /bookings [post]
func (h *BookingHandler) CreateBooking(c *gin.Context) {
	var req BookingRequest
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

// UpdateBookingStatus godoc
// @Summary      Update booking status
// @Description  Admin/Mod cập nhật trạng thái booking
// @Tags         bookings
// @Accept       json
// @Produce      json
// @Param        id path int true "Booking ID"
// @Param        request body UpdateBookingStatusRequest true "New status"
// @Success      200 {object} map[string]interface{} "booking updated"
// @Failure      400 {object} map[string]string "invalid input"
// @Failure      500 {object} map[string]string "internal server error"
// @Router       /bookings/{id}/status [put]
func (h *BookingHandler) UpdateBookingStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid.booking_id"})
		return
	}

	var req UpdateBookingStatusRequest
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

// GetBookingByID godoc
// @Summary      Get booking by ID
// @Description  Lấy chi tiết 1 booking
// @Tags         bookings
// @Produce      json
// @Param        id path int true "Booking ID"
// @Success      200 {object} map[string]interface{} "booking detail"
// @Failure      400 {object} map[string]string "invalid id"
// @Failure      404 {object} map[string]string "not found"
// @Router       /bookings/{id} [get]
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

// GetBookingByUserID godoc
// @Summary      Get bookings by user
// @Description  Lấy tất cả booking của user hiện tại
// @Tags         bookings
// @Produce      json
// @Success      200 {object} map[string]interface{} "list of bookings"
// @Failure      401 {object} map[string]string "unauthorized"
// @Router       /bookings/me [get]
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

// GetAllBooking godoc
// @Summary      Get all bookings
// @Description  Admin/Moderator lấy tất cả booking
// @Tags         bookings
// @Produce      json
// @Success      200 {object} map[string]interface{} "list of all bookings"
// @Failure      403 {object} map[string]string "forbidden"
// @Failure      500 {object} map[string]string "internal server error"
// @Router       /bookings [get]
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

func (h *BookingHandler) CheckAvailability(c *gin.Context) {
	var req dto.CheckAvailabilityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request data"})
		return
	}

	 // Parse string -> time.Time
    start, err := time.Parse(time.RFC3339, req.StartTime)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": "invalid start_time"})
        return
    }
    end, err := time.Parse(time.RFC3339, req.EndTime)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": "invalid end_time"})
        return
    }
	unavailable, err := h.usecase.CheckAvailability(c.Request.Context(), req.SpaceIDs, start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, dto.CheckAvailabilityResponse{
		UnavailableSpaceIDs: unavailable,
	})
}