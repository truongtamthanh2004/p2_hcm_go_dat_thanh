package route

import (
	"map-service/internal/handler"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(mapHandler *handler.MapHandler) *gin.Engine {
	r := gin.Default()
	r.GET("/api/v1/map/venues", mapHandler.ListVenues)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return r
}
