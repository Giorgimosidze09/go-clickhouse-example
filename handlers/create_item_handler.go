package handlers

import (
	"go-clickhouse-example/models"
	"go-clickhouse-example/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// @Security BearerAuth
// CreateItem godoc
// @Summary Create a new item
// @Description Create a new item and save it to the database, then publish to NATS
// @Tags items
// @Accept  json
// @Produce  json
// @Param item body models.ItemRequest true "Item to create"
// @Success 201 {object} models.ItemResponse "Created item"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /items [post]
func (h *ItemHandler) CreateItem(c *gin.Context) {
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
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to create an item"})
		return
	}

	// Bind the incoming request to the item model
	var itemRequest models.ItemRequest
	if err := c.ShouldBindJSON(&itemRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	item := models.ItemResponse{
		Name:  itemRequest.Name,
		Price: itemRequest.Price,
	}

	// Save item to database
	err = h.DBService.SaveItem(&item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save item to database"})
		return
	}

	// Publish the item to NATS
	err = h.NATSService.PublishItem(item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish item to NATS"})
		return
	}

	// Return the created item as a response
	c.JSON(http.StatusCreated, item)
}
