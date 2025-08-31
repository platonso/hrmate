package main

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/platonso/hrmate/internal/domain"
	"log"
	"net/http"
)

func (app *application) HandleUpdateUserStatus(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, errors.New("invalid UUID"))
		return
	}

	user, err := app.repositories.Users.FindById(userID)
	if err != nil {
		WriteJSONError(w, http.StatusNotFound, domain.ErrUserNotFound)
		return
	}

	var statusUserDTO UserStatusUpdateRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&statusUserDTO); err != nil {
		WriteJSONError(w, http.StatusBadRequest, errors.New("invalid JSON"))
		return
	}

	user.UpdateStatus(statusUserDTO.Status)

	if err := app.repositories.Users.Update(user); err != nil {
		WriteJSONError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

func (app *application) HandleGetUsers(w http.ResponseWriter, r *http.Request) {
	role, ok := GetUserRole(r.Context())
	if !ok {
		WriteJSONError(w, http.StatusUnauthorized, errors.New("missing role"))
		return
	}

	var rolesToQuery []domain.Role
	switch role {
	case domain.RoleAdmin:
		rolesToQuery = []domain.Role{domain.RoleHR, domain.RoleEmployee}
	case domain.RoleHR:
		rolesToQuery = []domain.Role{domain.RoleEmployee}
	default:
		WriteJSONError(w, http.StatusForbidden, errors.New("role not authorized"))
		return
	}

	users, err := app.repositories.Users.FindAllByRole(rolesToQuery...)
	if err != nil {
		WriteJSONError(w, http.StatusInternalServerError, err)
		log.Printf("failed to find users by role: %v", err)
		return
	}

	if users == nil {
		users = []domain.User{}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(users); err != nil {
		log.Printf("failed to encode response: %v", err)
	}

}
