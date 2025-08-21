package route

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"notification-service/config"
	"notification-service/internal/handler"
	"strconv"
)

func SetupRouter(notificationHandler *handler.NotificationHandler) *gin.Engine {
	router := gin.Default()
	router.HandleMethodNotAllowed = true // return 405 on wrong method

	router.POST("/api/v1/notifications", notificationHandler.SendNotification)
	router.GET("/api/v1/notifications/:userId", notificationHandler.GetNotifications)

	// WebSocket connect
	router.GET("/ws/:userId", func(c *gin.Context) {
		userIDStr := c.Param("userId")
		userID, err := strconv.Atoi(userIDStr)
		if err != nil || userID <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid userId"})
			return
		}
		config.HandleWS(c.Writer, c.Request, uint(userID))
	})

	return router
}
