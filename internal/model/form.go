package model

import (
	"github.com/google/uuid"
	"time"
)

type StatusType string

const (
	Pending  StatusType = "pending"
	Approved StatusType = "approved"
)

type Form struct {
	ID          uuid.UUID  `db:"id" json:"id"`
	UserID      uuid.UUID  `db:"user_id" json:"userId"`
	Title       string     `db:"title" json:"title"`
	Description string     `db:"description" json:"description"`
	CreatedAt   time.Time  `db:"created_at" json:"createdAt"`
	Status      StatusType `db:"status" json:"status"`
}

func NewForm(userId uuid.UUID, title, description string) (Form, error) {
	return Form{
		ID:          uuid.New(),
		UserID:      userId,
		Title:       title,
		Description: description,
		CreatedAt:   time.Now(),
		Status:      Pending,
	}, nil
}

func (f *Form) Approve() {
	f.Status = Approved
}
