package handler

import (
	"context"
	"net/http"
	"strconv"
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
