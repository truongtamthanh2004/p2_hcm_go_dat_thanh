package proxy

import (
	"github.com/gin-gonic/gin"
	"net/http/httputil"
	"net/url"
)

func NewReverseProxy(target string) gin.HandlerFunc {
	return func(c *gin.Context) {
		remote, err := url.Parse(target)
		if err != nil {
			c.JSON(500, gin.H{"error": "Invalid service URL"})
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(remote)

		// Update request to forward path properly
		path := c.Param("path")
		if path == "" {
			path = "/"
		} else if path[0] != '/' {
			path = "/" + path
		}
		c.Request.URL.Path = path

		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
