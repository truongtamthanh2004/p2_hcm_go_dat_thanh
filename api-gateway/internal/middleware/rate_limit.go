// middleware/rate_limit.go
package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type clientData struct {
	requests int
	resetAt  time.Time
}

var (
	clients = make(map[string]*clientData)
	mu      sync.Mutex
)

// RateLimitMiddleware creates a middleware that limits requests per IP per minute
func RateLimitMiddleware(maxRequests int) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		now := time.Now()

		mu.Lock()
		data, exists := clients[ip]
		if !exists || now.After(data.resetAt) {
			// New client or time window expired → reset counter
			clients[ip] = &clientData{
				requests: 1,
				resetAt:  now.Add(time.Minute),
			}
			mu.Unlock()
			c.Next()
			return
		}

		if data.requests >= maxRequests {
			mu.Unlock()
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":   "Too many requests",
				"message": "Rate limit exceeded. Please wait before trying again.",
			})
			return
		}

		data.requests++
		mu.Unlock()

		c.Next()
	}
}
