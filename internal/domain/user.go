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

type EmployeeRepository interface {
	Create(user *User) error
	FindById(userId uuid.UUID) (*User, error)
	FindByEmail(email string) (*User, error)
	FindAdmin(admin Role) (*User, error)
	Update(user *User) error
	FindAllByRole(...Role) ([]User, error)
	IsActive(userID uuid.UUID) (bool, error)
}

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

func NewEmployee(firstName, lastName, position, email, password string) User {
	return User{
		ID:             uuid.New(),
		Role:           RoleEmployee,
		FirstName:      firstName,
		LastName:       lastName,
		Position:       position,
		Email:          email,
		HashedPassword: password,
		IsActive:       true,
	}
}

func NewHR(firstName, lastName, position, email, password string) User {
	return User{
		ID:             uuid.New(),
		Role:           RoleHR,
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

func (u *User) UpdateStatus(isActive bool) {
	u.IsActive = isActive
}
