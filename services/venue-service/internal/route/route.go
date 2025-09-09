package route

import (
	"venue-service/internal/handler"
	"venue-service/internal/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)


func SetupRouter(venueHandler *handler.VenueHandler, spaceHandler *handler.SpaceHandler, amenityHandler *handler.AmenityHandler) *gin.Engine {
	r := gin.Default()
	v := r.Group("/api/v1/venues")
	{
		v.POST("", middleware.RequireAuth("user"), venueHandler.CreateVenue)
		v.GET("", middleware.RequireAuth("user"), venueHandler.GetVenues)
		v.GET("/:id", middleware.RequireAuth("user"), venueHandler.GetVenueByID)
		v.PUT("/:id", middleware.RequireAuth("user"), venueHandler.UpdateVenue)
		v.DELETE("/:id", middleware.RequireAuth("user"), venueHandler.DeleteVenue)

		// Amenities in venue
		v.POST("/:id/amenities", middleware.RequireAuth("user"), venueHandler.AddAmenity)
		v.DELETE("/:id/amenities/:venueAmenityId", middleware.RequireAuth("user"), venueHandler.RemoveAmenity)

		// Spaces under venue
		v.POST("/:id/spaces", middleware.RequireAuth("user"), spaceHandler.CreateSpace)
	}

	s := r.Group("/api/v1/spaces")
	{
		s.GET("/:id", spaceHandler.GetSpace)
		s.GET("/search", spaceHandler.SearchSpaces)
		s.PUT("/:id", middleware.RequireAuth("user"), spaceHandler.UpdateSpace)
		s.DELETE("/:id", middleware.RequireAuth("user"), spaceHandler.DeleteSpace)

		// manager update
		s.PUT("/:id/manager", middleware.RequireAuth("user"), spaceHandler.UpdateManager)
	}

	//admin
	a := r.Group("/api/v1/admin/amenities")
	{
		a.POST("", middleware.RequireAuth("admin", "moderator"), amenityHandler.CreateAmenity)
		a.GET("", middleware.RequireAuth("admin", "moderator"), amenityHandler.GetAllAmenities)
		a.GET("/:id", middleware.RequireAuth("admin", "moderator"), amenityHandler.GetAmenity)
		a.PUT("/:id", middleware.RequireAuth("admin", "moderator"), amenityHandler.UpdateAmenity)
		a.DELETE("/:id", middleware.RequireAuth("admin", "moderator"), amenityHandler.DeleteAmenity)
	}

	admin := r.Group("/api/v1/admin/venues")
	{
		// GET /admin/venues?status=pending
		admin.GET("", middleware.RequireAuth("admin", "moderator"), venueHandler.ListVenues)

		// PUT /admin/venues/:id/approve
		admin.PUT("/:id/approve", middleware.RequireAuth("admin", "moderator"), venueHandler.ApproveVenue)

		// PUT /admin/venues/:id/block
		admin.PUT("/:id/block", middleware.RequireAuth("admin", "moderator"), venueHandler.BlockVenue)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return r
}
