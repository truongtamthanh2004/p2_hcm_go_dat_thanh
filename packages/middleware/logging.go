package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggingMiddleware(c *gin.Context) {
	startTime := time.Now()

	c.Next()

	endTime := time.Now()
	latency := endTime.Sub(startTime)

	statusCode := c.Writer.Status()
	method := c.Request.Method
	path := c.Request.URL.Path
	clientIP := c.ClientIP()

	log.Printf("[GIN] %s | %3d | %13v | %15s | %-7s %s\n",
		endTime.Format("2006/01/02 - 15:04:05"),
		statusCode,
		latency,
		clientIP,
		method,
		path,
	)
}
