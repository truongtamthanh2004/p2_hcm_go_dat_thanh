package handler

import (
	"net/http"
	"strconv"
	"venue-service/internal/usecase"

	"github.com/gin-gonic/gin"
)

type SpaceHandler struct {
	uc usecase.SpaceUsecase
}

func NewSpaceHandler(uc usecase.SpaceUsecase) *SpaceHandler {
	return &SpaceHandler{uc: uc}
}

func (h *SpaceHandler) GetSpaceByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid.space_id"})
		return
	}

	space, err := h.uc.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "space.not_found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": space, "message": "space.success"})
}
