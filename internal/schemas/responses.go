package schemas

import "github.com/google/uuid"

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"user with such email already exists"`
}

type LoginResponse struct {
	Token        string `json:"token" example:"eyJhbGciOiJI..."`
	RefreshToken string `json:"refreshToken" example:"eyJhbGciOiJI"`
}
