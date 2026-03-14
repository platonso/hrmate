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
	ID          uuid.UUID `json:"id"          db:"id"`
	UserID      uuid.UUID `json:"userId"      db:"user_id"`
	Title       string    `json:"title"       db:"title"`
	Description string    `json:"description" db:"description"`

	StartDate *time.Time `json:"startDate"   db:"start_date"`
	EndDate   *time.Time `json:"endDate"     db:"end_date"`

	CreatedAt  time.Time  `json:"createdAt"   db:"created_at"`
	ApprovedAt *time.Time `json:"approvedAt"  db:"approved_at"`
	Status     FormStatus `json:"status"      db:"status"`
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

func (f *Form) UpdateStatus(newStatus FormStatus) {
	approveTime := time.Now()
	f.Status = newStatus
	f.ApprovedAt = &approveTime
}
