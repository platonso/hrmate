package model

import (
	"fmt"
	"github.com/google/uuid"
)

type RoleType string

const (
	RoleEmployee RoleType = "employee"
	RoleHR       RoleType = "hr"
)

func IsValidRole(r RoleType) bool {
	return r == RoleEmployee || r == RoleHR
}

type User struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Role      RoleType  `db:"role" json:"role"`
	FirstName string    `db:"first_name" json:"firstName"`
	LastName  string    `db:"last_name" json:"lastName"`
	Position  string    `db:"position" json:"position"`
	Email     string    `db:"email" json:"email"`
	Password  string    `db:"password" json:"-"`
}

func NewUser(role RoleType, firstName, lastName, position, email, password string) (User, error) {
	if !IsValidRole(role) {
		return User{}, fmt.Errorf("invalid role: %s", role)
	}

	return User{
		ID:        uuid.New(),
		Role:      role,
		FirstName: firstName,
		LastName:  lastName,
		Position:  position,
		Email:     email,
		Password:  password,
	}, nil
}

func (u *User) ChangeFirstName(newFirstName string) {
	u.FirstName = newFirstName
}

func (u *User) ChangeLastName(newLastName string) {
	u.LastName = newLastName
}

func (u *User) ChangePosition(newPosition string) {
	u.Position = newPosition
}
