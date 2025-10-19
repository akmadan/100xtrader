package dto

import (
	"time"
)

// UserCreateRequest represents the request to create a new user
type UserCreateRequest struct {
	Name  string  `json:"name" validate:"required,min=2,max=255"`
	Email string  `json:"email" validate:"required,email"`
	Phone *string `json:"phone,omitempty" validate:"omitempty,min=10,max=20"`
}

// UserUpdateRequest represents the request to update a user
type UserUpdateRequest struct {
	Name  *string `json:"name,omitempty" validate:"omitempty,min=2,max=255"`
	Email *string `json:"email,omitempty" validate:"omitempty,email"`
	Phone *string `json:"phone,omitempty" validate:"omitempty,min=10,max=20"`
}

// UserResponse represents the response for user data
type UserResponse struct {
	ID           int        `json:"id"`
	Name         string     `json:"name"`
	Email        string     `json:"email"`
	Phone        *string    `json:"phone,omitempty"`
	LastSignedIn *time.Time `json:"last_signed_in,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}

// UserListResponse represents the response for listing users
type UserListResponse struct {
	Users []UserResponse `json:"users"`
	Total int            `json:"total"`
}

// UserSignInRequest represents the request for user sign in
type UserSignInRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// UserSignInResponse represents the response for user sign in
type UserSignInResponse struct {
	User       UserResponse `json:"user"`
	LastSignIn time.Time    `json:"last_sign_in"`
	IsNewUser  bool         `json:"is_new_user"`
}
