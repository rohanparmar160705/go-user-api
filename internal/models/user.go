/*
Package models defines the data structures (DTOs) used for API requests and responses.
It includes struct tags for:
- JSON marshaling/unmarshaling
- Input validation using the "validator" package (e.g., required fields, date formats)
*/
package models

// CreateUserRequest represents the request body for creating a user
type CreateUserRequest struct {
	Name string `json:"name" validate:"required,min=2,max=100"`
	DOB  string `json:"dob" validate:"required"`
}

// UpdateUserRequest represents the request body for updating a user
type UpdateUserRequest struct {
	Name string `json:"name" validate:"required,min=2,max=100"`
	DOB  string `json:"dob" validate:"required"`
}

// UserResponse represents the response for a single user
type UserResponse struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
	DOB  string `json:"dob"`
	Age  *int   `json:"age,omitempty"` // Optional, only for GET requests
}

// UsersListResponse represents the response for listing users with pagination
type UsersListResponse struct {
	Data       []UserResponse `json:"data"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
	TotalPages int            `json:"total_pages"`
}
