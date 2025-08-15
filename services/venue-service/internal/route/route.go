package route

import (
	"github.com/gin-gonic/gin"
	"p2_hcm_go_dat_thanh/services/venue-service/internal/handler"
)

func SetupRouter(venueHandler *handler.VenueHandler) *gin.Engine {
	router := gin.Default()
	router.HandleMethodNotAllowed = true // return 405 on wrong method

	api := router.Group("/api/v1/venues")
	{
		// CREATE
		api.POST("/", venueHandler.CreateVenue)

		// READ (single)
		api.GET("/:id", venueHandler.GetVenue)

		// UPDATE
		api.PUT("/:id", venueHandler.UpdateVenue)

		// DELETE
		api.DELETE("/:id", venueHandler.DeleteVenue)

		// LIST / SEARCH
		api.GET("/", venueHandler.SearchVenues)
	}

	return router
}
