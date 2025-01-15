package handlers

import (
	"go-clickhouse-example/utils"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// @Security BearerAuth
// DeleteItem godoc
// @Summary Delete an item
// @Description Remove an item from the database by its ID
// @Tags items
// @Produce json
// @Param id path string true "Item ID"
// @Success 200 {object} map[string]string "Item deleted successfully"
// @Failure 400 {object} map[string]string "Invalid item ID"
// @Failure 404 {object} map[string]string "Item not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /items/{id} [delete]
func (h *ItemHandler) DeleteItem(c *gin.Context) {
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
	// Parse and validate the JWT token
	user, err := utils.ParseJWT(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	// Check user role for access control (optional)
	if user.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to delete this item"})
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

	// Delete the item from the database
	err = h.DBService.DeleteItem(itemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete item"})
		return
	}

	// Publish the item deletion to NATS
	err = h.NATSService.PublishItem(item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish item deletion to NATS"})
		return
	}

	// Return success message
	c.JSON(http.StatusOK, gin.H{"message": "Item deleted successfully"})
}
