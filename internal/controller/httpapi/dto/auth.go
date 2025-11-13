package dto

import "github.com/platonso/hrmate/internal/domain"

type UserStatusUpdateRequest struct {
	Status bool `json:"status"`
}

type RegisterRequest struct {
	FirstName string      `json:"firstName" validate:"required,min=2"`
	LastName  string      `json:"lastName" validate:"required,min=2"`
	Position  string      `json:"position" validate:"required,min=2"`
	Email     string      `json:"email" validate:"required,email"`
	Password  string      `json:"password" validate:"required"`
	Role      domain.Role `json:"role" validate:"required,oneof=employee hr"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	Token string `json:"token"`
}
