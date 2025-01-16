package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"go-clickhouse-example/models"

	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
)

// SearchItems handles the search, filter, and sorting functionality with pagination
// @Summary Search, filter, and sort items with pagination
// @Description Searches, filters, sorts, and paginates items based on query parameters
// @Tags Items
// @Accept json
// @Produce json
// @Param search query string false "Search query"
// @Param min_price query float64 false "Minimum price"
// @Param max_price query float64 false "Maximum price"
// @Param sort_by query string false "Sort by field (e.g., price, name)"
// @Param sort_order query string false "Sort order (ASC or DESC)"
// @Param page query int false "Page number (default is 1)"
// @Param limit query int false "Items per page (default is 10)"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "List of items with pagination metadata"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /items/search [get]
func (h *ItemHandler) SearchItems(c *gin.Context) {
	// Get search parameters from the query string
	searchQuery := c.DefaultQuery("search", "")
	minPrice := c.DefaultQuery("min_price", "0")
	maxPrice := c.DefaultQuery("max_price", "100000")
	sortBy := c.DefaultQuery("sort_by", "price")
	sortOrder := c.DefaultQuery("sort_order", "ASC")

	// Validate sortOrder to prevent SQL injection
	if sortOrder != "ASC" && sortOrder != "DESC" {
		sortOrder = "ASC" // Default to ASC if invalid
	}

	// Validate sortBy to ensure it's a valid column
	validSortBy := map[string]bool{
		"id":    true,
		"name":  true,
		"price": true,
	}
	if !validSortBy[sortBy] {
		sortBy = "price" // Default to "price" if invalid
	}

	// Build the SQL query dynamically
	query := "SELECT id, name, price FROM items WHERE 1=1"
	params := []interface{}{}

	// Add full-text search conditions if a search query is provided
	if searchQuery != "" {
		query += " AND (positionCaseInsensitive(name, ?) > 0)"
		params = append(params, searchQuery)
	}

	// Add price filtering
	query += " AND price BETWEEN ? AND ?"
	params = append(params, minPrice, maxPrice)

	// Add sorting
	query += fmt.Sprintf(" ORDER BY %s %s", sortBy, sortOrder)

	// Execute the query
	rows, err := h.DBService.Query(query, params...)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	defer rows.Close()

	// Parse the result and return the items
	var items []models.ItemResponse
	for rows.Next() {
		var item models.ItemResponse
		if err := rows.Scan(&item.ID, &item.Name, &item.Price); err != nil {
			log.Printf("Error scanning row: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
		items = append(items, item)
	}

	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *ItemHandler) PublishItemSearchResults(c *gin.Context) error {
	// Connect to NATS
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return err
	}
	defer nc.Close()

	var items []models.ItemResponse

	// Publish items to a subject
	for _, item := range items {
		data, err := json.Marshal(item)
		if err != nil {
			return err
		}

		// Publish each item to NATS (can be a batch operation as well)
		if err := nc.Publish("item.search.results", data); err != nil {
			return err
		}
	}

	return nil
}
