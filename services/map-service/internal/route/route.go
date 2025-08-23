package route

import (
	"map-service/internal/handler"

	"github.com/gin-gonic/gin"
)

func SetupRouter(mapHandler *handler.MapHandler) *gin.Engine {
	r := gin.Default()
	r.GET("/api/v1/map/venues", mapHandler.ListVenues)
	return r
}
