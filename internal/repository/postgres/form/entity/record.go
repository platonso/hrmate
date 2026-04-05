package entity

import (
	"time"

	"github.com/google/uuid"
)

type FormRecord struct {
	ID          uuid.UUID `db:"id"`
	UserID      uuid.UUID `db:"user_id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`

	StartDate *time.Time `db:"start_date"`
	EndDate   *time.Time `db:"end_date"`

	CreatedAt  time.Time  `db:"created_at"`
	ReviewedAt *time.Time `db:"reviewed_at"`
	Status     string     `db:"status"`
	Comment    *string    `db:"comment"`
}
