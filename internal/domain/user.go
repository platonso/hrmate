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
	ID             uuid.UUID
	Role           Role
	FirstName      string
	LastName       string
	Position       string
	Email          string
	HashedPassword string
	IsActive       bool
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
		user.IsActive = true
	}

	return user
}

func (u *User) ChangeNames(newFirstName, newLastName string) {
	u.FirstName = newFirstName
	u.LastName = newLastName
}

func (u *User) Activate() bool {
	if u.IsActive {
		return false
	}

	u.IsActive = true
	return true
}

func (u *User) Deactivate() bool {
	if !u.IsActive {
		return false
	}

	u.IsActive = false
	return true
}
