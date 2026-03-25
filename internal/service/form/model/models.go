package model

import (
	"time"

	"github.com/platonso/hrmate/internal/domain"
)

type FormCreateInput struct {
	Title       string
	Description string
	StartDate   *time.Time
	EndDate     *time.Time
}

type FormWithUser struct {
	User domain.User
	Form domain.Form
}

type FormsWithUser struct {
	User  domain.User
	Forms []domain.Form
}
