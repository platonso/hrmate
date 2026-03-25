package entity

import (
	"github.com/google/uuid"
)

type UserRecord struct {
	ID             uuid.UUID `db:"id"`
	Role           string    `db:"user_role"`
	FirstName      string    `db:"first_name"`
	LastName       string    `db:"last_name"`
	Position       string    `db:"position"`
	Email          string    `db:"email"`
	HashedPassword string    `db:"hashed_password"`
	IsActive       bool      `db:"is_active"`
}
