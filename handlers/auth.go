// handlers/auth.go
package handlers

import (
	"fmt"
	"go-clickhouse-example/models"
	"go-clickhouse-example/services"
	"go-clickhouse-example/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication-related requests
type AuthHandler struct {
	AuthService *services.AuthService
}

// NewAuthHandler creates a new AuthHandler instance
func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{AuthService: authService}
}

// RegisterUser godoc
// @Summary Register a new user
// @Description Registers a new user with a role and returns a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.UserRequest true "User to register"
// @Success 201 {string} string "JWT Token"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /register [post]
func (h *AuthHandler) RegisterUser(c *gin.Context) {
	var userRequest models.UserRequest
	if err := c.ShouldBindJSON(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Hash the password
	hashedPassword, err := utils.HashPassword(userRequest.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create a user model for storage
	user := &models.User{
		Username: userRequest.Username,
		Password: hashedPassword,
		Role:     userRequest.Role,
	}

	// Register user using the AuthService
	createdUser, err := h.AuthService.RegisterUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(createdUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user":  createdUser,
		"token": token,
	})
}

// LoginUser godoc
// @Summary Login user and get JWT token
// @Description Logs in the user and returns a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.UserRequest true "User login credentials"
// @Success 200 {string} string "JWT Token"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /login [post]
func (h *AuthHandler) LoginUser(c *gin.Context) {
	var userRequest models.UserRequest
	if err := c.ShouldBindJSON(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	fmt.Println("Login request:", userRequest) // Add this line to see the parsed request

	// Authenticate user
	user, err := h.AuthService.AuthenticateUser(userRequest.Username, userRequest.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
