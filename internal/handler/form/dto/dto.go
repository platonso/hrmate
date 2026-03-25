package dto

import (
	"time"

	"github.com/platonso/hrmate/internal/domain"
)

type FormCreateRequest struct {
	Title       string     `json:"title" validate:"required,min=3"`
	Description string     `json:"description" validate:"required"`
	StartDate   *time.Time `json:"startDate"`
	EndDate     *time.Time `json:"endDate"`
}

type FormCommentRequest struct {
	Comment string `json:"comment" validate:"required"`
}

//type FormResponse struct {
//	ID          uuid.UUID `json:"id"`
//	UserID      uuid.UUID `json:"userId"`
//	Title       string    `json:"title"`
//	Description string    `json:"description"`
//
//	StartDate *time.Time `json:"startDate"`
//	EndDate   *time.Time `json:"endDate"`
//
//	CreatedAt  time.Time  `json:"createdAt"`
//	ApprovedAt *time.Time `json:"approvedAt"`
//	Status     string     `json:"status"`
//	Comment    *string    `json:"comment"`
//}

type FormsWithUserResponse struct {
	User  domain.User   `json:"user"`
	Forms []domain.Form `json:"forms"`
}
