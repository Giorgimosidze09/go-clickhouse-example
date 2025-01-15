package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"go-clickhouse-example/utils"

	"github.com/gin-gonic/gin"
)

// @Security BearerAuth
// GetItem godoc
// @Summary Get an item by ID
// @Description Retrieve a single item from the database by its ID
// @Tags items
// @Produce  json
// @Param id path string true "Item ID"
// @Success 200 {object} models.ItemResponse "Retrieved item"
// @Failure 400 {object} map[string]string "Invalid item ID"
// @Failure 404 {object} map[string]string "Item not found"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /items/{id} [get]
func (h *ItemHandler) GetItem(c *gin.Context) {
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

	// Now the tokenString is the actual JWT token
	fmt.Println("Extracted token:", tokenString)

	// Parse and validate the JWT token
	user, err := utils.ParseJWT(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	// Get the item ID from the path
	id := c.Param("id")
	itemID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	// Retrieve the item from the database
	item, err := h.DBService.GetItemByID(itemID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	// Optionally, check user role for access control
	// Example: Only admin can access all items
	if user.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to access this item"})
		return
	}

	// Publish the item to NATS (if needed)
	err = h.NATSService.PublishItem(item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish item to NATS"})
		return
	}

	// Return the item as a response
	c.JSON(http.StatusOK, item)
}
