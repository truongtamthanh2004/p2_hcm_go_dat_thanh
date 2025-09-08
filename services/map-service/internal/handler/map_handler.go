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

// ListVenues godoc
// @Summary      List Venues
// @Description  Get all venues with their location info
// @Tags         venues
// @Produce      json
// @Success      200 {object} map[string]interface{} "venues list"
// @Failure      500 {object} map[string]string "internal server error"
// @Router       /venues [get]
func (h *MapHandler) ListVenues(c *gin.Context) {
	venues, err := h.Usecase.GetVenuesWithLocation()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "venues.found", "data": venues})
}
