package middleware

import (
	"go-clickhouse-example/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware is used to protect routes that require authentication
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from the Authorization header
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			c.Abort()
			return
		}

		// Remove the "Bearer " prefix if it exists
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		// Parse and validate the JWT token
		user, err := utils.ParseJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token in middleware"})
			c.Abort()
			return
		}

		// Set the user info in the context for use in other handlers
		c.Set("user_id", user.ID)
		c.Set("role", user.Role)

		// Continue to the next handler
		c.Next()
	}
}
