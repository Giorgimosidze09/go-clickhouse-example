// middleware/rbac.go
package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RBACMiddleware is used to restrict access based on user roles
func RBACMiddleware(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user role from context
		role := c.MustGet("role").(string)

		// Check if the user's role is allowed to access this route
		for _, allowedRole := range roles {
			if role == allowedRole {
				c.Next()
				return
			}
		}

		// If the user's role is not allowed, return forbidden error
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to access this resource"})
		c.Abort()
	}
}
