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
	ID             uuid.UUID `json:"id"`
	Role           Role      `json:"role,omitempty"`
	FirstName      string    `json:"firstName"`
	LastName       string    `json:"lastName"`
	Position       string    `json:"position"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"-"`
	IsActive       bool      `json:"isActive"`
}

func NewUser(role Role, firstName, lastName, position, email, password string) User {
	user := User{
		ID:             uuid.New(),
		Role:           role,
		FirstName:      firstName,
		LastName:       lastName,
		Position:       position,
		Email:          email,
		HashedPassword: password,
		IsActive:       false,
	}

	if user.Role == RoleEmployee {
		user.ChangeStatus(true)
	}

	return user
}

func (u *User) ChangeNames(newFirstName, newLastName string) {
	u.FirstName = newFirstName
	u.LastName = newLastName
}

func (u *User) ChangeStatus(isActive bool) {
	u.IsActive = isActive
}
