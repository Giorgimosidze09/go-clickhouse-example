package services

import (
	"errors"
	"fmt"
	"go-clickhouse-example/models"
	"go-clickhouse-example/utils"
)

// AuthService handles authentication-related operations
type AuthService struct {
	DBService *DBService
}

// NewAuthService creates a new AuthService instance
func NewAuthService(dbService *DBService) *AuthService {
	return &AuthService{DBService: dbService}
}

// RegisterUser handles user registration and saves user to the database
func (s *AuthService) RegisterUser(user *models.User) (*models.UserResponse, error) {
	// Hash the password before saving to the database
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	user.Password = hashedPassword

	// Save user to the database
	err = s.DBService.SaveUser(user)
	if err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}

	// Return a response object, the user ID is now automatically generated
	return &models.UserResponse{
		ID:       user.ID, // This will be automatically set after insertion
		Username: user.Username,
		Role:     user.Role,
	}, nil
}

// LoginUser handles user login by checking the password
func (s *AuthService) LoginUser(userRequest models.UserRequest) (*models.UserResponse, error) {
	// Get the user from the database by username
	user, err := s.DBService.GetUserByUsername(userRequest.Username)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Check if the entered password matches the hashed password in the database
	if !utils.CheckPasswordHash(userRequest.Password, user.Password) {
		return nil, errors.New("invalid password")
	}

	// Return the user object (you may also want to return a JWT token here)
	return &user, nil
}

// AuthenticateUser authenticates a user and returns a JWT token
func (s *AuthService) AuthenticateUser(username, password string) (*models.UserResponse, error) {
	user, err := s.DBService.GetUserByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	fmt.Println("Entered password:", password)
	fmt.Println("Stored hashed password:", user.Password)
	if !utils.CheckPasswordHash(password, user.Password) {
		fmt.Println("Password mismatch!")
	}

	return &models.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Role:     user.Role,
	}, nil
}
