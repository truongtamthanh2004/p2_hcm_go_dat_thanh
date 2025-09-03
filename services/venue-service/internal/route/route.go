package route

import (
	"venue-service/internal/handler"

	"github.com/gin-gonic/gin"
)

func SetupRouter(venueHandler *handler.VenueHandler, spaceHandler *handler.SpaceHandler) *gin.Engine {
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

	router.GET("/api/v1/spaces/:id", spaceHandler.GetSpaceByID)

	return router
}
