package middleware

import (
	"net/http"
	"strings"
	"venue-service/internal/utils"

	"github.com/gin-gonic/gin"
)

func RequireAuth(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "error.missing_token"})
			c.Abort()
			return
		}
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "error.invalid_token"})
			c.Abort()
			return
		}
		tokenStr := tokenParts[1]

		claims, err := utils.ValidateToken(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "error.invalid_token"})
			c.Abort()
			return
		}

		if !claims.IsVerified {
			c.JSON(http.StatusForbidden, gin.H{"message": "error.user_account_is_not_verified"})
			c.Abort()
			return
		}

		if !claims.IsActive {
			c.JSON(http.StatusForbidden, gin.H{"message": "error.user_is_not_activated"})
			c.Abort()
			return
		}

		if len(allowedRoles) > 0 {
			allowed := false
			for _, role := range allowedRoles {
				if claims.Role == role {
					allowed = true
					break
				}
			}
			if !allowed {
				c.JSON(http.StatusForbidden, gin.H{"message": "error.forbidden"})
				c.Abort()
				return
			}
		}

		c.Set("userEmail", claims.Email)
		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}
