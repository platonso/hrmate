package dto

import "github.com/google/uuid"

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Role      string    `json:"role"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Position  string    `json:"position"`
	Email     string    `json:"email"`
	IsActive  bool      `json:"isActive"`
}
