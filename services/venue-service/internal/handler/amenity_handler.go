package handler

import (
	"net/http"
	"strconv"
	"venue-service/internal/constant"
	"venue-service/internal/dto"
	"venue-service/internal/usecase"

	"github.com/gin-gonic/gin"
)

type AmenityHandler struct {
	uc usecase.AmenityUsecase
}

func NewAmenityHandler(uc usecase.AmenityUsecase) *AmenityHandler {
	return &AmenityHandler{uc}
}

// @Summary Create a new amenity
// @Description Create a new amenity (Admin only)
// @Tags Admin Amenity
// @Accept json
// @Produce json
// @Param body body dto.CreateAmenityRequest true "CreateAmenityRequest"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/amenities [post]
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

// @Summary Get all amenities
// @Description Get list of all amenities (Admin/Moderator only)
// @Tags Admin Amenity
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/amenities [get]
func (h *AmenityHandler) GetAllAmenities(c *gin.Context) {
	amenities, err := h.uc.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": amenities})
}

// @Summary Get amenity by ID
// @Description Get detailed information of an amenity by ID (Admin/Moderator only)
// @Tags Admin Amenity
// @Produce json
// @Param id path int true "Amenity ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /admin/amenities/{id} [get]
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

// @Summary Update an amenity
// @Description Update an existing amenity by ID (Admin/Moderator only)
// @Tags Admin Amenity
// @Accept json
// @Produce json
// @Param id path int true "Amenity ID"
// @Param body body dto.UpdateAmenityRequest true "UpdateAmenityRequest"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/amenities/{id} [put]
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

// @Summary Delete an amenity
// @Description Delete an amenity by ID (Admin/Moderator only)
// @Tags Admin Amenity
// @Produce json
// @Param id path int true "Amenity ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/amenities/{id} [delete]
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
