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
