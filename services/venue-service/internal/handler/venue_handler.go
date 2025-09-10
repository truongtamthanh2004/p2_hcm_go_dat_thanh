package handler

import (
	"net/http"
	"strconv"
	"venue-service/internal/constant"
	"venue-service/internal/dto"
	"venue-service/internal/usecase"

	"github.com/gin-gonic/gin"
)

type VenueHandler struct {
	uc usecase.VenueUsecase
}

func NewVenueHandler(uc usecase.VenueUsecase) *VenueHandler {
	return &VenueHandler{uc}
}

// @Summary Create a new venue
// @Description Create a new venue (user must be authenticated)
// @Tags Venue
// @Accept json
// @Produce json
// @Param body body dto.CreateVenueRequest true "CreateVenueRequest"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /venues [post]
func (h *VenueHandler) CreateVenue(c *gin.Context) {
	var req dto.CreateVenueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": constant.ErrBadRequest.Error()})
		return
	}
	userID, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": constant.ErrUnauthorized.Error()})
		return
	}
	venue, err := h.uc.Create(c.Request.Context(), userID.(uint), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "venue created",
		"data":    venue,
	})
}

// @Summary Get all venues
// @Description Get list of venues accessible by the user
// @Tags Venue
// @Produce json
// @Param city query string false "Filter by city"
// @Param name query string false "Filter by name"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /venues [get]
func (h *VenueHandler) GetVenues(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": constant.ErrUnauthorized.Error()})
		return
	}
	city := c.Query("city")
	name := c.Query("name")

	venues, err := h.uc.GetAll(c.Request.Context(), userID.(uint), city, name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, venues)
}

// @Summary Get venue by ID
// @Description Get detailed info of a venue by ID
// @Tags Venue
// @Produce json
// @Param id path int true "Venue ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /venues/{id} [get]
func (h *VenueHandler) GetVenueByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": constant.ErrInvalidID.Error()})
		return
	}
	venue, err := h.uc.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, venue)
}

// @Summary Update a venue
// @Description Update a venue by ID (user must be authenticated)
// @Tags Venue
// @Accept json
// @Produce json
// @Param id path int true "Venue ID"
// @Param body body dto.UpdateVenueRequest true "UpdateVenueRequest"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /venues/{id} [put]
func (h *VenueHandler) UpdateVenue(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": constant.ErrInvalidID.Error()})
		return
	}
	var req dto.UpdateVenueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": constant.ErrBadRequest.Error()})
		return
	}
	userID, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": constant.ErrUnauthorized.Error()})
		return
	}
	venue, err := h.uc.Update(c.Request.Context(), userID.(uint), uint(id), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "venue updated",
		"data":    venue,
	})
}

// @Summary Delete a venue
// @Description Delete a venue by ID (user must be authenticated)
// @Tags Venue
// @Produce json
// @Param id path int true "Venue ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /venues/{id} [delete]
func (h *VenueHandler) DeleteVenue(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": constant.ErrInvalidID.Error()})
		return
	}
	userID, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": constant.ErrUnauthorized.Error()})
		return
	}
	if err := h.uc.Delete(c.Request.Context(), userID.(uint), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "venue deleted"})
}

// @Summary Add an amenity to a venue
// @Description Add an existing amenity to a venue (user must be authenticated)
// @Tags Venue Amenity
// @Accept json
// @Produce json
// @Param id path int true "Venue ID"
// @Param body body dto.AddAmenityRequest true "AddAmenityRequest"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /venues/{id}/amenities [post]
func (h *VenueHandler) AddAmenity(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": constant.ErrInvalidID.Error()})
		return
	}
	var req dto.AddAmenityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": constant.ErrBadRequest.Error()})
		return
	}
	userID, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": constant.ErrUnauthorized.Error()})
		return
	}
	if err := h.uc.AddAmenity(c.Request.Context(), userID.(uint), uint(id), req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "amenity added"})
}

// @Summary Remove an amenity from a venue
// @Description Remove a specific amenity from a venue (user must be authenticated)
// @Tags Venue Amenity
// @Produce json
// @Param venueAmenityId path int true "Venue Amenity ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /venues/{id}/amenities/{venueAmenityId} [delete]
func (h *VenueHandler) RemoveAmenity(c *gin.Context) {
	venueAmenityID, err := strconv.Atoi(c.Param("venueAmenityId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": constant.ErrInvalidID.Error()})
		return
	}
	userID, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": constant.ErrUnauthorized.Error()})
		return
	}
	if err := h.uc.RemoveAmenity(c.Request.Context(), userID.(uint), uint(venueAmenityID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "amenity removed"})
}

// @Summary List venues with filter
// @Description List venues by status (Admin/Moderator only)
// @Tags Admin Venue
// @Produce json
// @Param status query string false "Status filter: pending / approved / blocked"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/venues [get]
func (h *VenueHandler) ListVenues(c *gin.Context) {
	var filter dto.FilterVenueRequest
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": constant.ErrBadRequest.Error()})
		return
	}

	venues, err := h.uc.List(c.Request.Context(), filter.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": venues})
}

// @Summary Approve a venue
// @Description Approve a venue by ID (Admin/Moderator only)
// @Tags Admin Venue
// @Produce json
// @Param id path int true "Venue ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/venues/{id}/approve [put]
func (h *VenueHandler) ApproveVenue(c *gin.Context) {
	venueID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": constant.ErrInvalidID.Error()})
		return
	}

	if err := h.uc.UpdateStatus(c.Request.Context(), uint(venueID), constant.APPROVED); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "venue approved"})
}

// @Summary Block a venue
// @Description Block a venue by ID (Admin/Moderator only)
// @Tags Admin Venue
// @Produce json
// @Param id path int true "Venue ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/venues/{id}/block [put]
func (h *VenueHandler) BlockVenue(c *gin.Context) {
	venueID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": constant.ErrInvalidID.Error()})
		return
	}

	if err := h.uc.UpdateStatus(c.Request.Context(), uint(venueID), constant.BLOCKED); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "venue blocked"})
}
