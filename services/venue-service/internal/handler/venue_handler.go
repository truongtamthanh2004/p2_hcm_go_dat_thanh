package handler

import (
	"errors"
	"net/http"
	"p2_hcm_go_dat_thanh/services/venue-service/internal/model"
	"p2_hcm_go_dat_thanh/services/venue-service/internal/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type VenueHandler struct {
	usecase usecase.VenueUsecase
}

func NewVenueHandler(u usecase.VenueUsecase) *VenueHandler {
	return &VenueHandler{usecase: u}
}

type CreateVenueRequest struct {
	Name        string `json:"name" binding:"required"`
	Address     string `json:"address" binding:"required"`
	City        string `json:"city"`
	Description string `json:"description"`
}

type UpdateVenueRequest struct {
	Name        string `json:"name" binding:"required"`
	Address     string `json:"address" binding:"required"`
	City        string `json:"city" binding:"required"`
	Status      string `json:"status" binding:"required"`
	Description string `json:"description"`
}

func (h *VenueHandler) CreateVenue(c *gin.Context) {
	//userID, exists := c.Get("user_id")
	//if !exists {
	//	c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	//	return
	//}

	var req CreateVenueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "data": []interface{}{}})
		return
	}

	venue := model.Venue{
		Name:        req.Name,
		Address:     req.Address,
		City:        req.City,
		Description: req.Description,
		//UserID:      userID.(uint),
		Status: "pending",
	}

	if err := h.usecase.CreateVenue(&venue); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "data": []interface{}{}})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "venue.created", "data": venue})
}

func (h *VenueHandler) GetVenue(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid.venue.id", "data": []interface{}{}})
		return
	}
	venue, err := h.usecase.GetVenue(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "invalid.data", "data": []interface{}{}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "venue.found", "data": venue})
}

func (h *VenueHandler) UpdateVenue(c *gin.Context) {
	// 1. parse id
	id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid.id", "data": []interface{}{}})
		return
	}
	id := uint(id64)

	// 2. bind request
	var req UpdateVenueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "data": []interface{}{}})
		return
	}

	// 3. fetch existing record
	existing, err := h.usecase.GetVenue(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "venue.not.found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "data": []interface{}{}})
		}
		return
	}

	// 4. map allowed fields only
	existing.Name = req.Name
	existing.Address = req.Address
	existing.City = req.City
	existing.Status = req.Status
	existing.Description = req.Description

	// 5. call the existing update service (keep service unchanged)
	if err := h.usecase.UpdateVenue(existing); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "data": []interface{}{}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "venue.updated", "data": existing})
}

func (h *VenueHandler) DeleteVenue(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid.venue.id", "data": []interface{}{}})
		return
	}
	
	if err := h.usecase.DeleteVenue(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "data": []interface{}{}})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *VenueHandler) SearchVenues(c *gin.Context) {
	city := c.Query("city")
	name := c.Query("name")
	venues, err := h.usecase.SearchVenues(city, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "data": []interface{}{}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "venues.found", "data": venues})
}
