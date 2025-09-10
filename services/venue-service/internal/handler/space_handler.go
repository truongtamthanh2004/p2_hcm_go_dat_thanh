package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"
	"venue-service/internal/constant"
	"venue-service/internal/dto"
	"venue-service/internal/usecase"

	"github.com/gin-gonic/gin"
)

type SpaceHandler struct {
	uc usecase.SpaceUsecase
}

func NewSpaceHandler(uc usecase.SpaceUsecase) *SpaceHandler {
	return &SpaceHandler{uc: uc}
}

// @Summary Create a new space under a venue
// @Description Create a new space under a specific venue (user must be authenticated)
// @Tags Space
// @Accept json
// @Produce json
// @Param id path int true "Venue ID"
// @Param body body dto.CreateSpaceRequest true "CreateSpaceRequest"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /venues/{id}/spaces [post]
func (h *SpaceHandler) CreateSpace(c *gin.Context) {
	var req dto.CreateSpaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": constant.ErrBadRequest.Error()})
		return
	}

	uid, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": constant.ErrUnauthorized.Error()})
		return
	}
	userID, ok := uid.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": constant.ErrUnauthorized.Error()})
		return
	}

	venueID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": constant.ErrInvalidID.Error()})
		return
	}

	ctx := context.Background()
	space, err := h.uc.Create(ctx, userID, uint(venueID), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "space created",
		"data":    space,
	})
}

// @Summary Get space by ID
// @Description Retrieve detailed information of a space by ID
// @Tags Space
// @Produce json
// @Param id path int true "Space ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /spaces/{id} [get]
func (h *SpaceHandler) GetSpace(c *gin.Context) {
	spaceID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": constant.ErrInvalidID.Error()})
		return
	}

	ctx := context.Background()
	space, err := h.uc.GetByID(ctx, uint(spaceID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": constant.ErrNotFound.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": space})
}

// @Summary Update a space
// @Description Update information of a space (user must be authenticated)
// @Tags Space
// @Accept json
// @Produce json
// @Param id path int true "Space ID"
// @Param body body dto.UpdateSpaceRequest true "UpdateSpaceRequest"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /spaces/{id} [put]
func (h *SpaceHandler) UpdateSpace(c *gin.Context) {
	var req dto.UpdateSpaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": constant.ErrBadRequest.Error()})
		return
	}

	uid, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": constant.ErrUnauthorized.Error()})
		return
	}
	userID, ok := uid.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": constant.ErrUnauthorized.Error()})
		return
	}

	spaceID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": constant.ErrInvalidID.Error()})
		return
	}

	ctx := context.Background()
	space, err := h.uc.Update(ctx, userID, uint(spaceID), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "space updated",
		"data":    space,
	})
}

// @Summary Delete a space
// @Description Delete a space by ID (user must be authenticated)
// @Tags Space
// @Produce json
// @Param id path int true "Space ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /spaces/{id} [delete]
func (h *SpaceHandler) DeleteSpace(c *gin.Context) {
	uid, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": constant.ErrUnauthorized.Error()})
		return
	}
	userID, ok := uid.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": constant.ErrUnauthorized.Error()})
		return
	}

	spaceID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": constant.ErrInvalidID.Error()})
		return
	}

	ctx := context.Background()
	if err := h.uc.Delete(ctx, userID, uint(spaceID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "space deleted"})
}

// @Summary Update manager of a space
// @Description Assign or update manager for a space (user must be authenticated)
// @Tags Space
// @Accept json
// @Produce json
// @Param id path int true "Space ID"
// @Param body body dto.UpdateManagerRequest true "UpdateManagerRequest"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /spaces/{id}/manager [put]
func (h *SpaceHandler) UpdateManager(c *gin.Context) {
	var req dto.UpdateManagerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": constant.ErrBadRequest.Error()})
		return
	}

	uid, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": constant.ErrUnauthorized.Error()})
		return
	}
	userID, ok := uid.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": constant.ErrUnauthorized.Error()})
		return
	}

	spaceID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": constant.ErrInvalidID.Error()})
		return
	}

	ctx := context.Background()
	if err := h.uc.UpdateManager(ctx, userID, uint(spaceID), req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "manager updated"})
}

// @Summary Search spaces
// @Description Search spaces with filters and availability in a time range
// @Tags Space
// @Produce json
// @Param name query string false "Space name"
// @Param city query string false "City"
// @Param address query string false "Address"
// @Param type query string false "Space type (private_office, meeting_room, desk)"
// @Param start_time query string true "Start time (RFC3339 format)"
// @Param end_time query string true "End time (RFC3339 format)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /spaces/search [get]
func (h *SpaceHandler) SearchSpaces(c *gin.Context) {
	name := c.Query("name")
	city := c.Query("city")
	address := c.Query("address")
	spaceType := c.Query("type")

	loc, _ := time.LoadLocation("Asia/Ho_Chi_Minh")

	startTime, err := time.Parse(time.RFC3339, c.Query("start_time"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid start_time"})
		return
	}
	startTime = startTime.In(loc)

	endTime, err := time.Parse(time.RFC3339, c.Query("end_time"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid end_time"})
		return
	}
	endTime = endTime.In(loc)

	if !startTime.Before(endTime) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "start_time must be before end_time",
		})
		return
	}
	spaces, err := h.uc.SearchSpaces(c.Request.Context(), name, city, address, spaceType, startTime, endTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    spaces,
	})
}
