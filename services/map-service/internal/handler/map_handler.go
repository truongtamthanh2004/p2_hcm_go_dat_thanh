package handler

import (
	"map-service/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MapHandler struct {
	Usecase *usecase.MapUsecase
}

func NewMapHandler(uc *usecase.MapUsecase) *MapHandler {
	return &MapHandler{Usecase: uc}
}

func (h *MapHandler) ListVenues(c *gin.Context) {
	venues, err := h.Usecase.GetVenuesWithLocation()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "venues.found", "data": venues})
}
