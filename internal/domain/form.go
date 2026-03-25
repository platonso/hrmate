package domain

import (
	"time"

	"github.com/google/uuid"
)

type FormStatus string

const (
	StatusPending  FormStatus = "pending"
	StatusApproved FormStatus = "approved"
)

type Form struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"userId"`
	Title       string    `json:"title"`
	Description string    `json:"description"`

	StartDate *time.Time `json:"startDate"`
	EndDate   *time.Time `json:"endDate"`

	CreatedAt  time.Time  `json:"createdAt"`
	ApprovedAt *time.Time `json:"approvedAt"`
	Status     FormStatus `json:"status"`
	Comment    *string    `json:"comment"`
}

func NewForm(userID uuid.UUID, title, description string, startDate, endDate *time.Time) Form {
	return Form{
		ID:          uuid.New(),
		UserID:      userID,
		Title:       title,
		Description: description,
		StartDate:   startDate,
		EndDate:     endDate,
		CreatedAt:   time.Now(),
		ApprovedAt:  nil,
		Status:      StatusPending,
	}
}

func (f *Form) ApproveForm(comment string) {
	approveTime := time.Now()
	f.ApprovedAt = &approveTime
	f.Status = StatusApproved
	f.Comment = &comment
}
