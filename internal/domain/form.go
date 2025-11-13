package domain

import (
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	Pending  Status = "Pending"
	Approved Status = "Approved"
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
	Status     Status     `json:"status"      db:"status"`
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
		Status:      Pending,
	}
}

func (f *Form) UpdateStatus(newStatus Status) {
	approveTime := time.Now()
	f.Status = newStatus
	f.ApprovedAt = &approveTime
}
