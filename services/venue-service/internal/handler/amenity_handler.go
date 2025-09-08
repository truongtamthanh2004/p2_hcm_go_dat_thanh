package handler

import (
	"venue-service/internal/constant"
	"venue-service/internal/dto"
	"venue-service/internal/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AmenityHandler struct {
	uc usecase.AmenityUsecase
}

func NewAmenityHandler(uc usecase.AmenityUsecase) *AmenityHandler {
	return &AmenityHandler{uc}
}

// POST /admin/amenities
func (h *AmenityHandler) CreateAmenity(c *gin.Context) {
	var req dto.CreateAmenityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": constant.ErrBadRequest.Error()})
		return
	}
	a, err := h.uc.Create(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "amenity created",
		"data":    a,
	})
}

// GET /amenities
func (h *AmenityHandler) GetAllAmenities(c *gin.Context) {
	amenities, err := h.uc.GetAll(c.Request.Context()) 
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": amenities})
}

// GET /amenities/:id
func (h *AmenityHandler) GetAmenity(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": constant.ErrInvalidID.Error()})
		return
	}
	a, err := h.uc.GetByID(c.Request.Context(), uint(id)) 
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": a})
}

// PUT /admin/amenities/:id
func (h *AmenityHandler) UpdateAmenity(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": constant.ErrInvalidID.Error()})
		return
	}
	var req dto.UpdateAmenityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": constant.ErrBadRequest.Error()})
		return
	}
	a, err := h.uc.Update(c.Request.Context(), uint(id), req) 
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "amenity updated", "data": a})
}

// DELETE /admin/amenities/:id
func (h *AmenityHandler) DeleteAmenity(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": constant.ErrInvalidID.Error()})
		return
	}
	if err := h.uc.Delete(c.Request.Context(), uint(id)); err != nil { // 
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "amenity deleted"})
}
