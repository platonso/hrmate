package main

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/platonso/hrmate/internal/domain"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
)

func (app *application) adminImplementation() error {
	admin, err := app.repositories.Users.FindAdmin(domain.RoleAdmin)
	if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
		return fmt.Errorf("failed to check admin: %w", err)
	}

	if admin != nil {
		return nil // Админ уже существует
	}

	email := os.Getenv("ADMIN_EMAIL")
	password := os.Getenv("ADMIN_PASSWORD")

	if email == "" || password == "" {
		return errors.New("ADMIN_EMAIL and ADMIN_PASSWORD must be set in environment")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash admin password: %w", err)
	}

	adminUser := domain.User{
		ID:             uuid.New(),
		Role:           domain.RoleAdmin,
		FirstName:      "Super",
		LastName:       "Admin",
		Position:       "Administrator",
		Email:          email,
		HashedPassword: string(hashedPassword),
		IsActive:       true,
	}

	if err := app.repositories.Users.Create(&adminUser); err != nil {
		return fmt.Errorf("failed to create admin: %w", err)
	}

	log.Println("Admin user created.")
	return nil
}
