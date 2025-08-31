package main

import (
	"encoding/json"
	"errors"
	"github.com/platonso/hrmate/internal/domain"
	"log"
	"net/http"
	"time"
)

type FormCreateRequestDTO struct {
	Title       string     `json:"title" validate:"required,min=3"`
	Description string     `json:"description"`
	StartDate   *time.Time `json:"startDate"`
	EndDate     *time.Time `json:"endDate"`
}

type FormStatusUpdateRequestDTO struct {
	Status domain.Status `json:"status"`
}

func (s *FormStatusUpdateRequestDTO) Validate() error {
	switch s.Status {
	case domain.Pending, domain.Approved:
		return nil
	default:
		return errors.New("invalid status: must be 'Pending' or 'Approved'")
	}
}

type UserStatusUpdateRequestDTO struct {
	Status bool `json:"status"`
}

type RegisterRequestDTO struct {
	FirstName string `json:"firstName" validate:"required,min=2"` // Создать глобальный валидатор (один на приложение)
	LastName  string `json:"lastName" validate:"required,min=2"`
	Position  string `json:"position" validate:"required,min=2"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required"`
}

type LoginRequestDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponseDTO struct {
	Token string `json:"token"`
}

type ErrResponseDTO struct {
	Message string    `json:"message"`
	Time    time.Time `json:"time"`
}

func WriteJSONError(w http.ResponseWriter, statusCode int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	errResponseDTO := ErrResponseDTO{
		Message: err.Error(),
		Time:    time.Now(),
	}

	if err := json.NewEncoder(w).Encode(errResponseDTO); err != nil {
		log.Printf("failed to encode error response: %v", err)
	}
}
