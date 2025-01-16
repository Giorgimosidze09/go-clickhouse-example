package handlers

import (
	"go-clickhouse-example/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// @Security BearerAuth
// GetItems godoc
// @Summary Get all items
// @Description Retrieve all items from the database
// @Tags items
// @Produce  json
// @Success 200 {array} models.ItemResponse "List of items"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /items [get]
func (h *ItemHandler) GetItems(c *gin.Context) {
	// Extract token from Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
		return
	}

	// Ensure the token is prefixed with "Bearer "
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
		return
	}

	// Parse and validate the JWT token
	user, err := utils.ParseJWT(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	// Optionally, check user role for access control (e.g., only admin can access all items)
	if user.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to access all items"})
		return
	}

	// Retrieve all items from the database
	items, err := h.DBService.GetAllItems()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch items"})
		return
	}

	// Return the list of items as a response
	c.JSON(http.StatusOK, items)
}
