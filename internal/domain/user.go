package domain

import (
	"github.com/google/uuid"
)

type Role string

const (
	RoleEmployee Role = "employee"
	RoleHR       Role = "hr"
	RoleAdmin    Role = "admin"
)

type User struct {
	ID             uuid.UUID `json:"id" db:"id"`
	Role           Role      `json:"role,omitempty" db:"user_role"`
	FirstName      string    `json:"firstName" db:"first_name"`
	LastName       string    `json:"lastName" db:"last_name"`
	Position       string    `json:"position" db:"position"`
	Email          string    `json:"email" db:"email"`
	HashedPassword string    `json:"-" db:"hashed_password"`
	IsActive       bool      `json:"isActive" db:"is_active"`
}

func NewUser(role Role, firstName, lastName, position, email, password string) User {
	return User{
		ID:             uuid.New(),
		Role:           role,
		FirstName:      firstName,
		LastName:       lastName,
		Position:       position,
		Email:          email,
		HashedPassword: password,
		IsActive:       false,
	}
}

func (u *User) ChangeNames(newFirstName, newLastName string) {
	u.FirstName = newFirstName
	u.LastName = newLastName
}

func (u *User) ChangeStatus(isActive bool) {
	u.IsActive = isActive
}
