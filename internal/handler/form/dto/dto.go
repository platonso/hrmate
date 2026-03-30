package dto

import (
	"time"

	"github.com/google/uuid"
)

type FormCreateRequest struct {
	Title       string     `json:"title" validate:"required,min=1"`
	Description string     `json:"description"`
	StartDate   *time.Time `json:"startDate"`
	EndDate     *time.Time `json:"endDate"`
}

type FormCommentRequest struct {
	Comment string `json:"comment"`
}

type FormResponse struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"userId"`
	Title       string    `json:"title"`
	Description string    `json:"description"`

	StartDate *time.Time `json:"startDate"`
	EndDate   *time.Time `json:"endDate"`

	CreatedAt  time.Time  `json:"createdAt"`
	ReviewedAt *time.Time `json:"reviewedAt"`
	Status     string     `json:"status"`
	Comment    *string    `json:"comment"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Role      string    `json:"role"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Position  string    `json:"position"`
	Email     string    `json:"email"`
	IsActive  bool      `json:"isActive"`
}

type FormsWithUserResponse struct {
	User  UserResponse   `json:"user"`
	Forms []FormResponse `json:"forms"`
}
