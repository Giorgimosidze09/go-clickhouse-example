// models/user.go
package models

// User represents the user entity in the system
type User struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// UserRequest is used for user registration (without password hashing)
type UserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// UserResponse is used for the response when fetching user data
type UserResponse struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
	Role     string `json:"role"`
}
