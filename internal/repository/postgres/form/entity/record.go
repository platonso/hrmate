package entity

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type FormRecord struct {
	ID          uuid.UUID `db:"id"`
	UserID      uuid.UUID `db:"user_id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`

	StartDate sql.NullTime `db:"start_date"`
	EndDate   sql.NullTime `db:"end_date"`

	CreatedAt  time.Time      `db:"created_at"`
	ApprovedAt sql.NullTime   `db:"approved_at"`
	Status     string         `db:"status"`
	Comment    sql.NullString `db:"comment"`
}
