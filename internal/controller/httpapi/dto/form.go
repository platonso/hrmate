package dto

import (
	"github.com/platonso/hrmate/internal/domain"
	"time"
)

type FormCreateRequest struct {
	Title       string     `json:"title" validate:"required,min=3"`
	Description string     `json:"description"`
	StartDate   *time.Time `json:"startDate"`
	EndDate     *time.Time `json:"endDate"`
}

type FormStatusUpdateRequest struct {
	Status domain.Status `json:"status" validate:"required,oneof=pending approved"`
}

type FormWithUserResponse struct {
	User domain.User `json:"user"`
	Form domain.Form `json:"form"`
}

type FormsWithUserResponse struct {
	User  domain.User   `json:"user"`
	Forms []domain.Form `json:"forms"`
}
