package handlers

import (
	"go-clickhouse-example/models"
	"go-clickhouse-example/utils"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// @Security BearerAuth
// UpdateItem godoc
// @Summary Update an existing item
// @Description Update an item in the database and publish the updated item to NATS
// @Tags items
// @Accept  json
// @Produce  json
// @Param id path string true "Item ID"
// @Param item body models.ItemRequest true "Updated item details"
// @Success 200 {object} map[string]string "Item updated successfully"
// @Failure 400 {object} map[string]string "Invalid input or item ID"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /items/{id} [put]
func (h *ItemHandler) UpdateItem(c *gin.Context) {
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
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to update this item"})
		return
	}

	// Get the item ID from the path
	id := c.Param("id")
	itemID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	// Bind the updated item details from the request body
	var item models.ItemResponse
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Update the item in the database
	err = h.DBService.UpdateItem(itemID, item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update item"})
		return
	}

	// Publish the updated item to NATS
	err = h.NATSService.PublishItem(item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish item to NATS"})
		return
	}

	// Return success message
	c.JSON(http.StatusOK, gin.H{"message": "Item updated successfully"})
}
